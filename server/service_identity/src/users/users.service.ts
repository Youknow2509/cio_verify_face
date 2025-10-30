import { Injectable, NotFoundException, ConflictException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../entities/user.entity';
import { Employee } from '../entities/employee.entity';
import { CreateUserDto } from '../dto/create-user.dto';
import { UpdateUserDto } from '../dto/update-user.dto';
import { RedisService } from '../redis/redis.service';
import { AuthService } from '../auth/auth.service';
import { v4 as uuidv4 } from 'uuid';

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    @InjectRepository(Employee)
    private employeeRepository: Repository<Employee>,
    private redisService: RedisService,
    private authService: AuthService,
  ) {}

  async create(createUserDto: CreateUserDto, companyId: string): Promise<User> {
    // Check if user with same email already exists
    const existingUser = await this.userRepository.findOne({
      where: { email: createUserDto.email },
    });
    if (existingUser) {
      throw new ConflictException('User with this email already exists');
    }

    // Hash password
    const salt = await this.authService.generateSalt();
    const passwordHash = await this.authService.hashPassword(createUserDto.password);

    // Create user
    const user = this.userRepository.create({
      email: createUserDto.email,
      salt,
      password_hash: passwordHash,
      full_name: createUserDto.full_name,
      role: createUserDto.role || 2, // Default to EMPLOYEE
    });
    const savedUser = await this.userRepository.save(user);

    // Create employee profile
    const employee = this.employeeRepository.create({
      user_id: savedUser.user_id,
      company_id: companyId,
      employee_code: createUserDto.employee_code || `EMP${Date.now()}`,
      department: createUserDto.department,
      position: createUserDto.position,
      phone: createUserDto.phone,
      hire_date: new Date(),
    });
    await this.employeeRepository.save(employee);

    // Clear cache
    await this.redisService.del(`company_users:${companyId}:*`);

    return savedUser;
  }

  async findAll(companyId: string, page: number = 1, limit: number = 10, search?: string): Promise<{
    data: User[];
    total: number;
    page: number;
    limit: number;
  }> {
    const cacheKey = `company_users:${companyId}:${page}:${limit}:${search || ''}`;
    const cached = await this.redisService.get(cacheKey);
    
    if (cached) {
      return JSON.parse(cached);
    }

    const queryBuilder = this.userRepository
      .createQueryBuilder('user')
      .leftJoinAndSelect('user.employee', 'employee')
      .where('employee.company_id = :companyId', { companyId });

    if (search) {
      queryBuilder.andWhere(
        'user.full_name ILIKE :search OR user.email ILIKE :search OR employee.employee_code ILIKE :search',
        { search: `%${search}%` }
      );
    }

    const [data, total] = await queryBuilder
      .orderBy('user.created_at', 'DESC')
      .skip((page - 1) * limit)
      .take(limit)
      .getManyAndCount();

    const result = {
      data,
      total,
      page,
      limit,
    };

    // Cache for 5 minutes
    await this.redisService.set(cacheKey, JSON.stringify(result), 300);

    return result;
  }

  async findOne(id: string): Promise<User> {
    const cacheKey = `user:${id}`;
    const cached = await this.redisService.get(cacheKey);
    
    if (cached) {
      return JSON.parse(cached);
    }

    const user = await this.userRepository.findOne({
      where: { user_id: id },
      relations: ['employee'],
    });

    if (!user) {
      throw new NotFoundException('User not found');
    }

    // Cache for 1 hour
    await this.redisService.set(cacheKey, JSON.stringify(user), 3600);

    return user;
  }

  async findByEmail(email: string): Promise<User> {
    return await this.userRepository.findOne({
      where: { email },
      relations: ['employee'],
    });
  }

  async findById(id: string): Promise<User> {
    return await this.userRepository.findOne({
      where: { user_id: id },
      relations: ['employee'],
    });
  }

  async update(id: string, updateUserDto: UpdateUserDto): Promise<User> {
    const user = await this.findOne(id);
    
    Object.assign(user, updateUserDto);
    const updatedUser = await this.userRepository.save(user);

    // Update employee profile if provided
    if (updateUserDto.department || updateUserDto.position || updateUserDto.phone) {
      const employee = await this.employeeRepository.findOne({
        where: { user_id: id },
      });
      if (employee) {
        Object.assign(employee, {
          department: updateUserDto.department,
          position: updateUserDto.position,
          phone: updateUserDto.phone,
        });
        await this.employeeRepository.save(employee);
      }
    }

    // Clear cache
    await this.redisService.del(`user:${id}`);
    await this.redisService.del(`company_users:${user.employee?.company_id}:*`);

    return updatedUser;
  }

  async remove(id: string): Promise<void> {
    const user = await this.findOne(id);
    
    // Soft delete by setting is_active to false
    user.is_active = false;
    await this.userRepository.save(user);

    // Clear cache
    await this.redisService.del(`user:${id}`);
    await this.redisService.del(`company_users:${user.employee?.company_id}:*`);
  }
}
