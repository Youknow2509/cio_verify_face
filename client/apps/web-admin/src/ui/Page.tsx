import Breadcrumb from 'react-bootstrap/Breadcrumb';
import Stack from 'react-bootstrap/Stack';
import type { ReactNode } from 'react';

interface PageProps {
  title: string;
  subtitle?: string;
  breadcrumb?: { label: string; path?: string }[];
  actions?: ReactNode;
  children: ReactNode;
}

export function Page({ title, subtitle, breadcrumb, actions, children }: PageProps) {
  return (
    <Stack gap={3} className="pb-4">
      <Stack direction="horizontal" className="justify-content-between align-items-start gap-3 flex-wrap">
        <div>
          {breadcrumb && breadcrumb.length > 0 && (
            <Breadcrumb className="mb-2">
              {breadcrumb.map((item, index) => (
                <Breadcrumb.Item
                  key={item.label}
                  linkAs="span"
                  active={index === breadcrumb.length - 1}
                >
                  {item.path && index !== breadcrumb.length - 1 ? (
                    <a href={item.path} className="text-decoration-none">
                      {item.label}
                    </a>
                  ) : (
                    item.label
                  )}
                </Breadcrumb.Item>
              ))}
            </Breadcrumb>
          )}
          <h1 className="fs-3 fw-semibold mb-1 text-primary">{title}</h1>
          {subtitle && <p className="text-secondary mb-0">{subtitle}</p>}
        </div>
        {actions && <div className="d-flex align-items-center gap-2 flex-wrap">{actions}</div>}
      </Stack>
      {children}
    </Stack>
  );
}
