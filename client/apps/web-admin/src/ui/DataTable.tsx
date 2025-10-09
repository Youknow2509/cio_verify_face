import Alert from 'react-bootstrap/Alert';
import Pagination from 'react-bootstrap/Pagination';
import Spinner from 'react-bootstrap/Spinner';
import Table from 'react-bootstrap/Table';
import Stack from 'react-bootstrap/Stack';
import type { ReactNode } from 'react';

export interface DataTableColumn<T> {
  header: ReactNode;
  accessor?: keyof T;
  render?: (row: T) => ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  columns: DataTableColumn<T>[];
  data: T[];
  loading?: boolean;
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  emptyMessage?: string;
  keySelector: (row: T) => string;
}

export function DataTable<T>({
  columns,
  data,
  loading = false,
  page,
  pageSize,
  total,
  onPageChange,
  emptyMessage = 'Không có dữ liệu',
  keySelector,
}: DataTableProps<T>) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  return (
    <Stack gap={3}>
      <div className="table-responsive">
        <Table striped hover responsive className="align-middle">
          <thead>
            <tr>
              {columns.map((col, index) => (
                <th key={index} className={col.className} scope="col">
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {loading && (
              <tr>
                <td colSpan={columns.length} className="text-center py-4">
                  <Spinner animation="border" size="sm" className="me-2" />
                  Đang tải dữ liệu...
                </td>
              </tr>
            )}
            {!loading && data.length === 0 && (
              <tr>
                <td colSpan={columns.length} className="p-0">
                  <Alert variant="light" className="mb-0 text-center py-4">
                    {emptyMessage}
                  </Alert>
                </td>
              </tr>
            )}
            {!loading &&
              data.map((row) => (
                <tr key={keySelector(row)}>
                  {columns.map((col, index) => (
                    <td key={index} className={col.className}>
                      {col.render
                        ? col.render(row)
                        : col.accessor
                        ? (row[col.accessor] as ReactNode)
                        : null}
                    </td>
                  ))}
                </tr>
              ))}
          </tbody>
        </Table>
      </div>

      {totalPages > 1 && (
        <div className="d-flex flex-wrap justify-content-between align-items-center gap-2">
          <span className="text-secondary small">
            Trang {page} / {totalPages} · {total} bản ghi
          </span>
          <Pagination className="mb-0">
            <Pagination.Prev disabled={page === 1} onClick={() => onPageChange(page - 1)} />
            {Array.from({ length: totalPages }).map((_, index) => (
              <Pagination.Item
                key={index}
                active={index + 1 === page}
                onClick={() => onPageChange(index + 1)}
              >
                {index + 1}
              </Pagination.Item>
            ))}
            <Pagination.Next
              disabled={page === totalPages}
              onClick={() => onPageChange(page + 1)}
            />
          </Pagination>
        </div>
      )}
    </Stack>
  );
}
