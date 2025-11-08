import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline, Box, Typography } from '@mui/material';

const theme = createTheme();

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ p: 4 }}>
        <Typography variant="h4">Employee Portal - Coming Soon</Typography>
        <Typography variant="body1">Cổng thông tin nhân viên</Typography>
      </Box>
    </ThemeProvider>
  );
}

export default App;
