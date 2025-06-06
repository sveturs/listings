// Менеджер для хранения файлов вне Redux store
class FileUploadManager {
  private files: Map<string, File> = new Map();

  addFile(fileId: string, file: File): void {
    this.files.set(fileId, file);
  }

  getFile(fileId: string): File | undefined {
    return this.files.get(fileId);
  }

  removeFile(fileId: string): void {
    this.files.delete(fileId);
  }

  clear(): void {
    this.files.clear();
  }
}

export const fileUploadManager = new FileUploadManager();
