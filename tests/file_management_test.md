# File Management System Test Guide

## 📁 **File Management System Implementation**

### **✅ Implemented Features**

1. **MongoDB GridFS Integration**
   - File storage in MongoDB GridFS
   - File metadata collection
   - Unique file naming system

2. **File Upload System**
   - Multipart form-data support
   - File validation (type & size)
   - Permission-based access control
   - Metadata storage

3. **File Management Operations**
   - Upload files to achievements
   - Get file lists by achievement
   - Download/stream files
   - Delete files with permission check

4. **Security & Validation**
   - File type validation: pdf, doc, docx, jpg, jpeg, png
   - File size limit: 5MB max
   - Role-based access control
   - User permission validation

### **🔧 API Endpoints**

#### **Upload File**
```http
POST /api/v1/achievements/{id}/files
Authorization: Bearer {mahasiswa_token}
Content-Type: multipart/form-data

Form Data:
- file: [binary file]
```

#### **Get Achievement Files**
```http
GET /api/v1/achievements/{id}/files
Authorization: Bearer {token}
```

#### **Download File**
```http
GET /api/v1/achievements/{id}/files/{fileId}/download
Authorization: Bearer {token}
```

#### **Delete File**
```http
DELETE /api/v1/achievements/{id}/files/{fileId}
Authorization: Bearer {mahasiswa_token}
```

### **🧪 Testing Instructions**

#### **1. Test File Upload**
```bash
# Using curl
curl -X POST http://localhost:8080/api/v1/achievements/{achievement_id}/files \
  -H "Authorization: Bearer {mahasiswa_token}" \
  -F "file=@certificate.pdf"

# Expected Response:
{
  "success": true,
  "message": "File uploaded successfully", 
  "data": {
    "id": "uuid",
    "filename": "certificate.pdf",
    "size": 1234567,
    "content_type": "application/pdf",
    "uploaded_at": "2025-11-28T10:00:00Z",
    "achievement_id": "achievement-uuid",
    "uploaded_by": "user-uuid",
    "gridfs_id": "gridfs-object-id"
  }
}
```

#### **2. Test File Validation**
```bash
# Test invalid file type
curl -X POST http://localhost:8080/api/v1/achievements/{achievement_id}/files \
  -H "Authorization: Bearer {mahasiswa_token}" \
  -F "file=@malicious.exe"

# Expected: 400 Error - invalid file type

# Test file size limit
curl -X POST http://localhost:8080/api/v1/achievements/{achievement_id}/files \
  -H "Authorization: Bearer {mahasiswa_token}" \
  -F "file=@large_file_6mb.pdf"

# Expected: 400 Error - file size exceeds 5MB limit
```

#### **3. Test Access Control**
```bash
# Mahasiswa trying to upload to other's achievement
curl -X POST http://localhost:8080/api/v1/achievements/{other_achievement}/files \
  -H "Authorization: Bearer {mahasiswa_token}" \
  -F "file=@certificate.pdf"

# Expected: 403 Error - Access denied
```

#### **4. Test File Download**
```bash
# Download file
curl -X GET http://localhost:8080/api/v1/achievements/{achievement_id}/files/{file_id}/download \
  -H "Authorization: Bearer {token}" \
  -o downloaded_file.pdf

# Should download the actual file
```

### **🔍 Database Verification**

#### **MongoDB Collections Created:**
1. **fs.files** - GridFS file metadata
2. **fs.chunks** - GridFS file chunks  
3. **achievement_files** - Custom file metadata

#### **Check MongoDB Data:**
```javascript
// Connect to MongoDB
use prestasi_files

// Check file metadata
db.achievement_files.find().pretty()

// Check GridFS files
db.fs.files.find().pretty()

// Count files by achievement
db.achievement_files.aggregate([
  { $group: { _id: "$achievement_id", count: { $sum: 1 } } }
])
```

### **🚨 Error Scenarios to Test**

1. **Missing file in request**: 400 - File is required
2. **Invalid achievement ID**: 404 - Achievement not found  
3. **Permission denied**: 403 - Access denied
4. **Invalid file type**: 400 - Invalid file type
5. **File too large**: 400 - File size exceeds limit
6. **MongoDB connection error**: 500 - Database error
7. **GridFS error**: 500 - File storage error

### **📊 Performance Considerations**

- **File Size Limit**: 5MB per file
- **GridFS Chunk Size**: 255KB (MongoDB default)
- **Concurrent Uploads**: Limited by server resources
- **Download Streaming**: Memory efficient for large files
- **File Type Validation**: Client-side + server-side

### **🔧 Troubleshooting**

#### **Common Issues:**
1. **MongoDB Connection**: Check MongoDB service status
2. **File Upload Fails**: Verify multipart form-data
3. **GridFS Errors**: Check MongoDB permissions
4. **Download Issues**: Verify file exists in GridFS

#### **Debug Commands:**
```bash
# Check MongoDB connection
docker ps | grep mongo

# Check file service logs
tail -f server.log | grep "FileService"

# Test MongoDB GridFS manually
mongosh prestasi_files
```

### **✅ System Status**
- ✅ File Upload Implementation: **Complete**
- ✅ File Download/Streaming: **Complete** 
- ✅ File Validation: **Complete**
- ✅ Access Control: **Complete**
- ✅ MongoDB GridFS Integration: **Complete**
- ✅ API Endpoints: **Complete**
- ✅ Error Handling: **Complete**

**File Management System is now fully operational! 🎉**