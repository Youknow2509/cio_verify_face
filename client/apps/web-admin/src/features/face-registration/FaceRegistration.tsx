// src/features/face-registration/FaceRegistration.tsx

import { useCallback, useEffect, useMemo, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Spinner from 'react-bootstrap/Spinner';
import { useNavigate, useSearchParams } from 'react-router-dom';

import styles from './FaceRegistration.module.scss';
import { Page } from '@/ui/Page';
import { useUi } from '@/app/providers/UiProvider';
import { getEmployee, getEmployeeFaceData, uploadFaceData } from '@/services';
import { clearAuthToken, HttpError } from '@/services/http';
import type { Employee, FaceData } from '@/types';

function FaceRegistration() {
  const [searchParams] = useSearchParams();
  const employeeId = searchParams.get('employeeId') ?? undefined;
  const navigate = useNavigate();
  const { showToast } = useUi();

  const [employee, setEmployee] = useState<Employee | null>(null);
  const [faces, setFaces] = useState<FaceData[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);

  const loadData = useCallback(async (id: string) => {
    try {
      setLoading(true);
      const res = await getEmployee(id);
      if (res.error) {
        showToast({ variant: 'warning', title: 'Không tìm thấy', message: 'Nhân viên không tồn tại' });
        navigate('/employees', { replace: true });
        return;
      }

      setEmployee(res.data);

      const facesRes = await getEmployeeFaceData(id);
      if (!facesRes.error && facesRes.data) {
        setFaces(facesRes.data);
      }
    } catch (err) {
      if (err instanceof HttpError && err.status === 401) {
        clearAuthToken();
        navigate('/login', { replace: true });
        return;
      }

      showToast({ variant: 'danger', title: 'Lỗi', message: 'Không thể tải dữ liệu' });
    } finally {
      setLoading(false);
    }
  }, [navigate, showToast]);

  useEffect(() => {
    if (!employeeId) {
      navigate('/employees', { replace: true });
      return;
    }

    void loadData(employeeId);
  }, [employeeId, loadData, navigate]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;
    setSelectedFiles(prev => [...prev, ...Array.from(files)]);
  };

  const handleRemoveSelected = (idx: number) => {
    setSelectedFiles(prev => prev.filter((_, i) => i !== idx));
  };

  const handleUpload = async () => {
    if (!employeeId || selectedFiles.length === 0) return;
    setUploading(true);
    try {
      for (const file of selectedFiles) {
        const res = await uploadFaceData(employeeId, file);
        if (!res.error && res.data) {
          setFaces(prev => [res.data!, ...prev]);
        }
      }

      showToast({ variant: 'success', title: 'Hoàn tất', message: 'Ảnh đã được tải lên' });
      // After successful upload, go back to employee detail
      navigate(`/employees/${employeeId}`);
    } catch (err) {
      if (err instanceof HttpError && err.status === 401) {
        clearAuthToken();
        navigate('/login', { replace: true });
        return;
      }

      showToast({ variant: 'danger', title: 'Lỗi', message: 'Không thể tải ảnh' });
    } finally {
      setUploading(false);
    }
  };

  const backToDetail = () => {
    if (employeeId) navigate(`/employees/${employeeId}`);
    else navigate('/employees');
  };

  const countText = useMemo(() => {
    const count = faces.length;
    return `Ảnh đã đăng ký (${count}/5)`;
  }, [faces]);

  return (
    <Page title="ĐĂNG KÝ DỮ LIỆU KHUÔN MẶT">
      <div className={styles.wrapper}>
        <button className={styles.backLink + ' btn btn-link p-0 mb-3'} onClick={backToDetail}>← Quay lại</button>

        {loading ? (
          <div className="d-flex align-items-center justify-content-center py-5">
            <Spinner animation="border" role="status" aria-hidden />
            <span className="ms-2">Đang tải...</span>
          </div>
        ) : (
          <>
            <section className={styles.employeeInfo}>
              <div className={styles.infoRow}>
                <div className={styles.avatar} aria-hidden>
                  <svg width="40" height="40" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 12a5 5 0 100-10 5 5 0 000 10zM3 20a9 9 0 0118 0H3z" />
                  </svg>
                </div>
                <div className={styles.infoText}>
                  <div className={styles.name}>{employee ? `${employee.name} (${employee.code})` : '---'}</div>
                  <div className={styles.meta}>{employee ? `${employee.email}  |  ${employee.department ?? ''}` : ''}</div>
                </div>
              </div>
            </section>

            <section className={styles.mainGrid}>
              <div className={styles.leftColumn}>
                <div className={styles.card}>
                  <h4>Upload ảnh mới</h4>
                  <div className={styles.uploadArea}>
                    <div className={styles.dropzone}>
                      <div className={styles.dropText}>Kéo thả ảnh vào đây</div>
                      <div className={styles.or}>hoặc</div>
                      <label className={styles.chooseBtn}>
                        <input type="file" accept="image/*" hidden multiple onChange={handleFileChange} />
                        Chọn ảnh
                      </label>
                    </div>
                    <div className={styles.supportText}>Hỗ trợ: JPG, PNG (Max: 5MB/ảnh)</div>

                    {selectedFiles.length > 0 && (
                      <div className="mt-3">
                        <div className="mb-2">Ảnh đã chọn:</div>
                        <div className="d-flex flex-wrap gap-2">
                          {selectedFiles.map((f, i) => (
                            <div key={i} className={styles.selectedThumb}>
                              <img src={URL.createObjectURL(f)} alt={f.name} className={styles.selectedThumbImg} />
                              <div className="d-flex gap-1 mt-1">
                                <small className="text-muted">{f.name}</small>
                                <button className="btn btn-link btn-sm p-0 ms-2" onClick={() => handleRemoveSelected(i)}>Xóa</button>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    <div className="mt-3">
                      <Button variant="outline-secondary" className="me-2" onClick={() => { /* TODO: webcam capture */ }}>
                        📷 Chụp bằng webcam
                      </Button>
                      <Button variant="primary" onClick={handleUpload} disabled={uploading || selectedFiles.length === 0}>
                        {uploading ? (
                          <span className="d-inline-flex align-items-center gap-2">
                            <Spinner animation="border" size="sm" role="status" aria-hidden />
                            Đang tải
                          </span>
                        ) : (
                          'Tải lên và quay lại'
                        )}
                      </Button>
                    </div>
                  </div>
                </div>
              </div>

              <div className={styles.rightColumn}>
                <div className={styles.card}>
                  <h4>Hướng dẫn chụp ảnh</h4>
                  <ul className={styles.guidelines}>
                    <li>Chụp thẳng, ánh sáng đủ</li>
                    <li>Không đeo kính, khẩu trang</li>
                    <li>Biểu cảm tự nhiên</li>
                    <li>5 góc độ khác nhau:
                      <ul>
                        <li>Thẳng</li>
                        <li>Nghiêng trái/phải</li>
                        <li>Hướng lên/xuống</li>
                      </ul>
                    </li>
                  </ul>
                </div>
              </div>
            </section>

            <section className={styles.gallerySection}>
              <div className={styles.count}>{countText}</div>

              <div className={styles.gallery}>
                {faces.length === 0 ? (
                  <div className="text-secondary">Chưa có ảnh nào.</div>
                ) : (
                  faces.map(face => (
                    <div className={styles.thumb} key={face.id}>
                      <div className={styles.thumbImg} style={{ backgroundImage: `url(${face.imageUrl})` }} />
                      <div className={styles.thumbMeta}>{face.fileName}</div>
                    </div>
                  ))
                )}
              </div>

              <div className={styles.qualityRow}>
                <label>Chất lượng nhận diện:</label>
                <div className={styles.qualityOptions}>
                  <label><input type="radio" name="quality" /> Trung bình</label>
                  <label><input type="radio" name="quality" /> Tốt</label>
                </div>
              </div>

            </section>
          </>
        )}
      </div>
    </Page>
  );
}

export default FaceRegistration;
