import React, { useState, useMemo } from 'react';
import {
    Box,
    Card,
    Typography,
    Chip,
    Grid,
    Stack,
    Tooltip,
    ToggleButton,
    ToggleButtonGroup,
    Button,
} from '@mui/material';
import { ViewWeek, ViewList } from '@mui/icons-material';
import type { Shift } from '@face-attendance/types';

interface ShiftTimelineProps {
    shifts: Shift[];
}

const DAYS_OF_WEEK_SHORT = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];
const HOURS = Array.from({ length: 24 }, (_, i) => i);

// Hàm để tìm range giờ làm việc
const getWorkingHoursRange = (
    shifts: Shift[]
): { start: number; end: number } => {
    if (shifts.length === 0) return { start: 8, end: 18 };

    let minHour = 23;
    let maxHour = 0;

    shifts.forEach((shift) => {
        const startHour = new Date(shift.start_time).getHours();
        const endHour = new Date(shift.end_time).getHours();
        minHour = Math.min(minHour, startHour);
        maxHour = Math.max(maxHour, endHour);
    });

    return {
        start: Math.max(0, minHour - 1),
        end: Math.min(24, maxHour + 1),
    };
};

const formatTime = (timeStr: string): string => {
    const date = new Date(timeStr);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
};

const getHourFromTime = (timeStr: string): number => {
    const date = new Date(timeStr);
    return date.getHours();
};

const getMinuteFromTime = (timeStr: string): number => {
    const date = new Date(timeStr);
    return date.getMinutes();
};

const SHIFT_COLORS = [
    '#FF6B6B',
    '#4ECDC4',
    '#45B7D1',
    '#FFA07A',
    '#98D8C8',
    '#F7DC6F',
    '#BB8FCE',
    '#85C1E2',
];

export const ShiftTimeline: React.FC<ShiftTimelineProps> = ({ shifts }) => {
    const [viewMode, setViewMode] = useState<'timeline' | 'table'>('timeline');

    // Nhóm ca theo ngày
    const shiftsByDay: { [key: number]: Shift[] } = {};
    for (let i = 1; i <= 7; i++) {
        shiftsByDay[i] = shifts.filter((shift) => shift.work_days.includes(i));
    }

    // Lấy range giờ làm việc để tối ưu hiển thị
    const workingHours = useMemo(() => getWorkingHoursRange(shifts), [shifts]);
    const displayHours = useMemo(
        () =>
            Array.from(
                { length: workingHours.end - workingHours.start },
                (_, i) => workingHours.start + i
            ),
        [workingHours]
    );

    // Tính chiều cao của ca dựa trên thời gian (40px per hour)
    const HOUR_HEIGHT = 40;
    const getShiftHeight = (shift: Shift): number => {
        const startHour = getHourFromTime(shift.start_time);
        const startMin = getMinuteFromTime(shift.start_time);
        const endHour = getHourFromTime(shift.end_time);
        const endMin = getMinuteFromTime(shift.end_time);

        const startTotal = startHour * 60 + startMin;
        const endTotal = endHour * 60 + endMin;
        const durationMins = endTotal - startTotal;

        return (durationMins / 60) * HOUR_HEIGHT;
    };

    const getShiftTop = (shift: Shift): number => {
        const startHour = getHourFromTime(shift.start_time);
        const startMin = getMinuteFromTime(shift.start_time);
        return ((startHour * 60 + startMin) / 60) * HOUR_HEIGHT;
    };

    return (
        <Card sx={{ mb: 3 }}>
            <Box p={1.5}>
                <Box
                    sx={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        mb: 1.5,
                    }}
                >
                    <Typography variant="subtitle1" fontWeight="bold">
                        Lịch ca làm việc tuần
                    </Typography>
                    <ToggleButtonGroup
                        value={viewMode}
                        exclusive
                        onChange={(_, newMode) =>
                            newMode && setViewMode(newMode)
                        }
                        size="small"
                        sx={{ bgcolor: '#f5f5f5' }}
                    >
                        <ToggleButton value="timeline" aria-label="timeline">
                            <ViewWeek fontSize="small" />
                        </ToggleButton>
                        <ToggleButton value="table" aria-label="table">
                            <ViewList fontSize="small" />
                        </ToggleButton>
                    </ToggleButtonGroup>
                </Box>

                {viewMode === 'timeline' ? (
                    // TIMELINE VIEW
                    <>
                        <Box sx={{ overflowX: 'auto' }}>
                            <Box sx={{ minWidth: 750 }}>
                                {/* Header với ngày */}
                                <Box sx={{ display: 'flex', mb: 0.5 }}>
                                    <Box
                                        sx={{
                                            minWidth: 60,
                                            fontSize: '12px',
                                            fontWeight: 'bold',
                                            py: 0.5,
                                        }}
                                    >
                                        Giờ
                                    </Box>
                                    {DAYS_OF_WEEK_SHORT.map((day) => (
                                        <Box
                                            key={day}
                                            sx={{
                                                flex: 1,
                                                minWidth: 100,
                                                textAlign: 'center',
                                                fontSize: '12px',
                                                fontWeight: '600',
                                                py: 0.5,
                                                color: 'primary.main',
                                            }}
                                        >
                                            {day}
                                        </Box>
                                    ))}
                                </Box>

                                {/* Timeline */}
                                <Box
                                    sx={{
                                        display: 'flex',
                                        border: '1px solid #e0e0e0',
                                        borderRadius: 1,
                                        overflow: 'hidden',
                                        bgcolor: '#fff',
                                    }}
                                >
                                    {/* Column giờ */}
                                    <Box
                                        sx={{
                                            minWidth: 60,
                                            borderRight: '2px solid #e0e0e0',
                                            bgcolor: '#f9f9f9',
                                        }}
                                    >
                                        {displayHours.map((hour) => (
                                            <Box
                                                key={hour}
                                                sx={{
                                                    height: HOUR_HEIGHT,
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    borderBottom:
                                                        '1px solid #e0e0e0',
                                                    fontSize: '10px',
                                                    fontWeight: '500',
                                                    color: '#666',
                                                }}
                                            >
                                                {`${hour
                                                    .toString()
                                                    .padStart(2, '0')}:00`}
                                            </Box>
                                        ))}
                                    </Box>

                                    {/* Columns ca theo ngày */}
                                    {DAYS_OF_WEEK_SHORT.map((_, dayIdx) => {
                                        const dayNum = dayIdx + 1;
                                        const dayCells = shiftsByDay[dayNum];

                                        return (
                                            <Box
                                                key={dayNum}
                                                sx={{
                                                    flex: 1,
                                                    minWidth: 100,
                                                    borderRight:
                                                        '1px solid #e0e0e0',
                                                    position: 'relative',
                                                }}
                                            >
                                                {/* Background giờ */}
                                                {displayHours.map((hour) => (
                                                    <Box
                                                        key={`bg-${hour}`}
                                                        sx={{
                                                            height: HOUR_HEIGHT,
                                                            borderBottom:
                                                                '1px solid #f0f0f0',
                                                            backgroundColor:
                                                                hour % 2 === 0
                                                                    ? '#fafafa'
                                                                    : '#fff',
                                                        }}
                                                    />
                                                ))}

                                                {/* Ca làm việc */}
                                                <Box
                                                    sx={{
                                                        position: 'absolute',
                                                        top: 0,
                                                        left: 0,
                                                        right: 0,
                                                        height:
                                                            displayHours.length *
                                                            HOUR_HEIGHT,
                                                    }}
                                                >
                                                    {dayCells.map(
                                                        (shift, idx) => {
                                                            const startHour =
                                                                new Date(
                                                                    shift.start_time
                                                                ).getHours();
                                                            const endHour =
                                                                new Date(
                                                                    shift.end_time
                                                                ).getHours();

                                                            // Skip nếu ca không nằm trong range hiển thị
                                                            if (
                                                                startHour >=
                                                                    workingHours.end ||
                                                                endHour <=
                                                                    workingHours.start
                                                            ) {
                                                                return null;
                                                            }

                                                            const height =
                                                                getShiftHeight(
                                                                    shift
                                                                );
                                                            const rawTop =
                                                                getShiftTop(
                                                                    shift
                                                                );
                                                            const offsetTop =
                                                                (startHour -
                                                                    workingHours.start) *
                                                                HOUR_HEIGHT;
                                                            const top =
                                                                offsetTop;
                                                            const bgColor =
                                                                SHIFT_COLORS[
                                                                    idx %
                                                                        SHIFT_COLORS.length
                                                                ];

                                                            return (
                                                                <Tooltip
                                                                    key={
                                                                        shift.shift_id
                                                                    }
                                                                    title={`${
                                                                        shift.name
                                                                    }\n${formatTime(
                                                                        shift.start_time
                                                                    )} - ${formatTime(
                                                                        shift.end_time
                                                                    )}`}
                                                                >
                                                                    <Box
                                                                        sx={{
                                                                            position:
                                                                                'absolute',
                                                                            top: `${top}px`,
                                                                            left: '1px',
                                                                            right: '1px',
                                                                            height: `${height}px`,
                                                                            backgroundColor:
                                                                                bgColor,
                                                                            border: `1px solid ${bgColor}`,
                                                                            borderRadius:
                                                                                '3px',
                                                                            padding:
                                                                                '2px 3px',
                                                                            cursor: 'pointer',
                                                                            opacity:
                                                                                shift.is_active
                                                                                    ? 1
                                                                                    : 0.5,
                                                                            transition:
                                                                                'all 0.2s ease',
                                                                            overflow:
                                                                                'hidden',
                                                                            '&:hover':
                                                                                {
                                                                                    boxShadow:
                                                                                        '0 2px 4px rgba(0,0,0,0.15)',
                                                                                    zIndex: 10,
                                                                                },
                                                                        }}
                                                                    >
                                                                        <Typography
                                                                            sx={{
                                                                                fontSize:
                                                                                    '9px',
                                                                                fontWeight:
                                                                                    'bold',
                                                                                color: '#fff',
                                                                                textShadow:
                                                                                    '0 1px 1px rgba(0,0,0,0.3)',
                                                                                display:
                                                                                    'block',
                                                                                whiteSpace:
                                                                                    'nowrap',
                                                                                overflow:
                                                                                    'hidden',
                                                                                textOverflow:
                                                                                    'ellipsis',
                                                                                lineHeight: 1.1,
                                                                            }}
                                                                        >
                                                                            {
                                                                                shift.name
                                                                            }
                                                                        </Typography>
                                                                        <Typography
                                                                            sx={{
                                                                                fontSize:
                                                                                    '7px',
                                                                                color: '#fff',
                                                                                textShadow:
                                                                                    '0 1px 1px rgba(0,0,0,0.3)',
                                                                                lineHeight: 1,
                                                                                display:
                                                                                    height >
                                                                                    28
                                                                                        ? 'block'
                                                                                        : 'none',
                                                                            }}
                                                                        >
                                                                            {formatTime(
                                                                                shift.start_time
                                                                            )}{' '}
                                                                            -{' '}
                                                                            {formatTime(
                                                                                shift.end_time
                                                                            )}
                                                                        </Typography>
                                                                    </Box>
                                                                </Tooltip>
                                                            );
                                                        }
                                                    )}
                                                </Box>
                                            </Box>
                                        );
                                    })}
                                </Box>
                            </Box>
                        </Box>
                    </>
                ) : (
                    // TABLE VIEW
                    <Box sx={{ overflowX: 'auto' }}>
                        <Box sx={{ minWidth: 500 }}>
                            <Box
                                sx={{
                                    display: 'grid',
                                    gridTemplateColumns: '80px 1fr 120px 120px',
                                    gap: 0,
                                    borderBottom: '2px solid #e0e0e0',
                                    bgcolor: '#f9f9f9',
                                }}
                            >
                                <Box
                                    sx={{
                                        padding: '8px',
                                        fontWeight: '600',
                                        fontSize: '12px',
                                        borderRight: '1px solid #e0e0e0',
                                    }}
                                >
                                    Ngày
                                </Box>
                                <Box
                                    sx={{
                                        padding: '8px',
                                        fontWeight: '600',
                                        fontSize: '12px',
                                        borderRight: '1px solid #e0e0e0',
                                    }}
                                >
                                    Ca làm việc
                                </Box>
                                <Box
                                    sx={{
                                        padding: '8px',
                                        fontWeight: '600',
                                        fontSize: '12px',
                                        textAlign: 'center',
                                        borderRight: '1px solid #e0e0e0',
                                    }}
                                >
                                    Giờ
                                </Box>
                                <Box
                                    sx={{
                                        padding: '8px',
                                        fontWeight: '600',
                                        fontSize: '12px',
                                        textAlign: 'center',
                                    }}
                                >
                                    Trạng thái
                                </Box>
                            </Box>
                            {DAYS_OF_WEEK_SHORT.map((day, dayIdx) => {
                                const dayNum = dayIdx + 1;
                                const dayCells = shiftsByDay[dayNum];

                                return (
                                    <Box key={dayNum}>
                                        {dayCells.length === 0 ? (
                                            <Box
                                                sx={{
                                                    display: 'grid',
                                                    gridTemplateColumns:
                                                        '80px 1fr 120px 120px',
                                                    gap: 0,
                                                    borderBottom:
                                                        '1px solid #e0e0e0',
                                                    bgcolor: '#fff',
                                                }}
                                            >
                                                <Box
                                                    sx={{
                                                        padding: '8px',
                                                        fontWeight: '600',
                                                        borderRight:
                                                            '1px solid #e0e0e0',
                                                        fontSize: '12px',
                                                    }}
                                                >
                                                    {day}
                                                </Box>
                                                <Box
                                                    sx={{
                                                        padding: '8px',
                                                        gridColumn: 'span 3',
                                                        color: '#999',
                                                        fontStyle: 'italic',
                                                        fontSize: '12px',
                                                    }}
                                                >
                                                    Không có ca
                                                </Box>
                                            </Box>
                                        ) : (
                                            dayCells.map((shift, idx) => (
                                                <Box
                                                    key={shift.shift_id}
                                                    sx={{
                                                        display: 'grid',
                                                        gridTemplateColumns:
                                                            '80px 1fr 120px 120px',
                                                        gap: 0,
                                                        borderBottom:
                                                            '1px solid #f0f0f0',
                                                        bgcolor:
                                                            idx % 2 === 0
                                                                ? '#fff'
                                                                : '#fafafa',
                                                    }}
                                                >
                                                    <Box
                                                        sx={{
                                                            padding: '8px',
                                                            borderRight:
                                                                '1px solid #e0e0e0',
                                                            fontSize: '12px',
                                                        }}
                                                    >
                                                        {idx === 0 && (
                                                            <strong>
                                                                {day}
                                                            </strong>
                                                        )}
                                                    </Box>
                                                    <Box
                                                        sx={{
                                                            padding: '8px',
                                                            borderRight:
                                                                '1px solid #e0e0e0',
                                                        }}
                                                    >
                                                        <Box
                                                            sx={{
                                                                display:
                                                                    'inline-block',
                                                                padding:
                                                                    '4px 8px',
                                                                backgroundColor:
                                                                    SHIFT_COLORS[
                                                                        idx %
                                                                            SHIFT_COLORS.length
                                                                    ],
                                                                color: '#fff',
                                                                borderRadius:
                                                                    '3px',
                                                                fontSize:
                                                                    '12px',
                                                                fontWeight:
                                                                    '500',
                                                            }}
                                                        >
                                                            {shift.name}
                                                        </Box>
                                                    </Box>
                                                    <Box
                                                        sx={{
                                                            padding: '8px',
                                                            textAlign: 'center',
                                                            borderRight:
                                                                '1px solid #e0e0e0',
                                                            fontSize: '12px',
                                                        }}
                                                    >
                                                        {formatTime(
                                                            shift.start_time
                                                        )}{' '}
                                                        -{' '}
                                                        {formatTime(
                                                            shift.end_time
                                                        )}
                                                    </Box>
                                                    <Box
                                                        sx={{
                                                            padding: '8px',
                                                            textAlign: 'center',
                                                        }}
                                                    >
                                                        <Box
                                                            sx={{
                                                                display:
                                                                    'inline-block',
                                                                padding:
                                                                    '2px 6px',
                                                                backgroundColor:
                                                                    shift.is_active
                                                                        ? '#c3e9cd'
                                                                        : '#f0f0f0',
                                                                color: shift.is_active
                                                                    ? '#1b5e20'
                                                                    : '#666',
                                                                borderRadius:
                                                                    '3px',
                                                                fontSize:
                                                                    '11px',
                                                                fontWeight:
                                                                    '500',
                                                            }}
                                                        >
                                                            {shift.is_active
                                                                ? 'Hoạt động'
                                                                : 'Tạm dừng'}
                                                        </Box>
                                                    </Box>
                                                </Box>
                                            ))
                                        )}
                                    </Box>
                                );
                            })}
                        </Box>
                    </Box>
                )}

                {/* Legend */}
                <Stack direction="row" spacing={1.5} sx={{ mt: 1.5 }}>
                    <Box display="flex" alignItems="center" gap={0.75}>
                        <Box
                            sx={{
                                width: 12,
                                height: 12,
                                backgroundColor: '#e3fcec',
                                border: '1px solid #51cf66',
                            }}
                        />
                        <Typography variant="caption" sx={{ fontSize: '11px' }}>
                            Ca hoạt động
                        </Typography>
                    </Box>
                    <Box display="flex" alignItems="center" gap={0.75}>
                        <Box
                            sx={{
                                width: 12,
                                height: 12,
                                backgroundColor: '#f5f5f5',
                                border: '1px solid #999',
                                opacity: 0.5,
                            }}
                        />
                        <Typography variant="caption" sx={{ fontSize: '11px' }}>
                            Ca tạm dừng
                        </Typography>
                    </Box>
                </Stack>
            </Box>
        </Card>
    );
};
