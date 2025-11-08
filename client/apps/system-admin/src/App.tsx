import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline, Box, Typography } from '@mui/material';

const theme = createTheme();

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ p: 4 }}>
        <Typography variant="h4">System Admin - Coming Soon</Typography>
        <Typography variant="body1">Giao diện quản trị hệ thống</Typography>
      </Box>
    </ThemeProvider>
  );
}

export default App;
