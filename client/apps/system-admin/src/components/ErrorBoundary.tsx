import { Component, ErrorInfo, ReactNode } from 'react';
import { Box, Typography, Button, Card, CardContent } from '@mui/material';
import { ErrorOutline, Refresh } from '@mui/icons-material';

interface Props {
    children: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { hasError: false, error: null };
    }

    static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('ErrorBoundary caught:', error, errorInfo);
    }

    handleReload = () => {
        window.location.reload();
    };

    handleGoHome = () => {
        window.location.href = '/dashboard';
    };

    render() {
        if (this.state.hasError) {
            return (
                <Box
                    sx={{
                        minHeight: '100vh',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        background: 'linear-gradient(135deg, #0f172a 0%, #1e1b4b 50%, #0f172a 100%)',
                        p: 3,
                    }}
                >
                    <Card sx={{ maxWidth: 500, textAlign: 'center' }}>
                        <CardContent sx={{ p: 4 }}>
                            <ErrorOutline sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
                            <Typography variant="h5" fontWeight="700" gutterBottom>
                                Đã xảy ra lỗi
                            </Typography>
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                                Ứng dụng gặp sự cố không mong muốn. Vui lòng thử tải lại trang hoặc quay về trang chủ.
                            </Typography>

                            {this.state.error && (
                                <Box
                                    sx={{
                                        bgcolor: 'action.hover',
                                        borderRadius: 1,
                                        p: 2,
                                        mb: 3,
                                        textAlign: 'left',
                                    }}
                                >
                                    <Typography variant="caption" color="error.main" component="pre" sx={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                                        {this.state.error.message}
                                    </Typography>
                                </Box>
                            )}

                            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
                                <Button
                                    variant="contained"
                                    startIcon={<Refresh />}
                                    onClick={this.handleReload}
                                >
                                    Tải lại trang
                                </Button>
                                <Button
                                    variant="outlined"
                                    onClick={this.handleGoHome}
                                >
                                    Về Dashboard
                                </Button>
                            </Box>
                        </CardContent>
                    </Card>
                </Box>
            );
        }

        return this.props.children;
    }
}
