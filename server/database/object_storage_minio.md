# Face Attendance System - Object Storage (Minio) Configuration
# File and media storage for face images, signatures, and reports

## Bucket Structure and Policies

### Bucket Organization
```bash
# Bucket structure for multi-tenant isolation
face-images/
├── {company_id}/
│   ├── {employee_id}/
│   │   ├── profile/
│   │   │   ├── original_{timestamp}.jpg
│   │   │   ├── processed_{timestamp}.jpg
│   │   │   └── thumbnail_{timestamp}.jpg
│   │   └── verification/
│   │       ├── {record_id}_{timestamp}.jpg
│   │       └── {record_id}_cropped_{timestamp}.jpg
│   └── temp/
│       └── upload_{session_id}.jpg

signatures/
├── {company_id}/
│   ├── {user_id}/
│   │   ├── signature_{timestamp}.png
│   │   └── signature_preview_{timestamp}.png
│   └── templates/
│       └── company_signature_template.png

reports/
├── {company_id}/
│   ├── daily/
│   │   └── attendance_report_{date}.pdf
│   ├── monthly/
│   │   └── summary_report_{year}_{month}.xlsx
│   ├── custom/
│   │   └── custom_report_{report_id}.pdf
│   └── exports/
│       └── data_export_{timestamp}.csv

backups/
├── database/
│   ├── daily/
│   │   └── pg_dump_{date}.sql.gz
│   ├── weekly/
│   │   └── full_backup_{date}.tar.gz
│   └── incremental/
│       └── wal_archive/
├── config/
│   └── system_config_{timestamp}.json
└── logs/
    └── audit_logs_{date}.json.gz

firmware/
├── devices/
│   ├── face_terminal_v1.2.3.bin
│   ├── face_terminal_v1.2.4.bin
│   └── changelog.json
└── updates/
    ├── pending/
    └── completed/
```

### Bucket Policies and Access Control
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "FaceImagesReadWrite",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::*:user/face-attendance-app"
      },
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::face-images/*"
    },
    {
      "Sid": "ReportsReadOnly",
      "Effect": "Allow", 
      "Principal": {
        "AWS": "arn:aws:iam::*:user/face-attendance-readonly"
      },
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::reports/*"
    },
    {
      "Sid": "BackupWriteOnly",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::*:user/face-attendance-backup"
      },
      "Action": [
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::backups/*"
    }
  ]
}
```

