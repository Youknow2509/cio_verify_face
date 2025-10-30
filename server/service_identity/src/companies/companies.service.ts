import { Injectable, NotFoundException, ConflictException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Company } from '../entities/company.entity';
import { CompanySecret } from '../entities/company-secret.entity';
import { CreateCompanyDto } from '../dto/create-company.dto';
import { UpdateCompanyDto } from '../dto/update-company.dto';
import { RedisService } from '../redis/redis.service';
import { v4 as uuidv4 } from 'uuid';
import * as bcrypt from 'bcrypt';

@Injectable()
export class CompaniesService {
  constructor(
    @InjectRepository(Company)
    private companyRepository: Repository<Company>,
    @InjectRepository(CompanySecret)
    private companySecretRepository: Repository<CompanySecret>,
    private redisService: RedisService,
  ) {}

  async create(createCompanyDto: CreateCompanyDto): Promise<Company> {
    // Check if company with same email already exists
    if (createCompanyDto.email) {
      const existingCompany = await this.companyRepository.findOne({
        where: { email: createCompanyDto.email },
      });
      if (existingCompany) {
        throw new ConflictException('Company with this email already exists');
      }
    }

    // Create company
    const company = this.companyRepository.create(createCompanyDto);
    const savedCompany = await this.companyRepository.save(company);

    // Create company secret
    const salt = await bcrypt.genSalt(10);
    const secret = uuidv4();
    const secretHash = await bcrypt.hash(secret, salt);

    const companySecret = this.companySecretRepository.create({
      company_id: savedCompany.company_id,
      salt,
      secret_hash: secretHash,
    });
    await this.companySecretRepository.save(companySecret);

    // Clear cache
    await this.redisService.del('companies_list:*');

    return savedCompany;
  }

  async findAll(page: number = 1, limit: number = 10, search?: string): Promise<{
    data: Company[];
    total: number;
    page: number;
    limit: number;
  }> {
    const cacheKey = `companies_list:${page}:${limit}:${search || ''}`;
    const cached = await this.redisService.get(cacheKey);
    
    if (cached) {
      return JSON.parse(cached);
    }

    const queryBuilder = this.companyRepository.createQueryBuilder('company');
    
    if (search) {
      queryBuilder.where(
        'company.name ILIKE :search OR company.email ILIKE :search',
        { search: `%${search}%` }
      );
    }

    const [data, total] = await queryBuilder
      .orderBy('company.created_at', 'DESC')
      .skip((page - 1) * limit)
      .take(limit)
      .getManyAndCount();

    const result = {
      data,
      total,
      page,
      limit,
    };

    // Cache for 2 minutes
    await this.redisService.set(cacheKey, JSON.stringify(result), 120);

    return result;
  }

  async findOne(id: string): Promise<Company> {
    const cacheKey = `company:${id}`;
    const cached = await this.redisService.get(cacheKey);
    
    if (cached) {
      return JSON.parse(cached);
    }

    const company = await this.companyRepository.findOne({
      where: { company_id: id },
      relations: ['employees', 'workShifts', 'devices'],
    });

    if (!company) {
      throw new NotFoundException('Company not found');
    }

    // Cache for 1 hour
    await this.redisService.set(cacheKey, JSON.stringify(company), 3600);

    return company;
  }

  async update(id: string, updateCompanyDto: UpdateCompanyDto): Promise<Company> {
    const company = await this.findOne(id);
    
    Object.assign(company, updateCompanyDto);
    const updatedCompany = await this.companyRepository.save(company);

    // Clear cache
    await this.redisService.del(`company:${id}`);
    await this.redisService.del('companies_list:*');

    return updatedCompany;
  }

  async remove(id: string): Promise<void> {
    const company = await this.findOne(id);
    
    // Soft delete by setting status to INACTIVE
    company.status = 'INACTIVE';
    await this.companyRepository.save(company);

    // Clear cache
    await this.redisService.del(`company:${id}`);
    await this.redisService.del('companies_list:*');
  }
}
