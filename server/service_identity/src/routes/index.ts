import { Router } from 'express';
import companiesRouter from './companies';
import usersRouter from './users';

const router = Router();

router.use('/companies', companiesRouter);
router.use('/users', usersRouter);

export default router;
