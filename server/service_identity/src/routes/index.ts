import { Router } from 'express';
import companiesRouter from './companies';
import usersRouter from './users';

const router = Router();

router.use('/companies', companiesRouter);
router.use('/users', usersRouter);
// Health check endpoint
router.get('/health', (req, res) => {
  res.status(200).send('OK');
});

export default router;
