'use client';

import { useState, useEffect } from 'react';
import {
  ChevronRightIcon,
  DocumentTextIcon,
  FolderIcon,
  FolderOpenIcon,
  Bars3Icon,
  XMarkIcon,
  ArrowLeftIcon,
} from '@heroicons/react/24/outline';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import AdminGuard from '@/components/AdminGuard';
import type { components } from '@/types/generated/api';

type DocFile = components['schemas']['internal_proj_docserver_handler.DocFile'];
type DocFilesResponse =
  components['schemas']['backend_pkg_utils.SuccessResponseSwag'] & {
    data?: components['schemas']['internal_proj_docserver_handler.DocFilesResponse'];
  };
type DocContentResponse =
  components['schemas']['backend_pkg_utils.SuccessResponseSwag'] & {
    data?: components['schemas']['internal_proj_docserver_handler.DocContentResponse'];
  };

export default function DocsPage() {
  const [files, setFiles] = useState<DocFile[]>([]);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [fileContent, setFileContent] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [expandedDirs, setExpandedDirs] = useState<Set<string>>(new Set());
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    fetchDocFiles();

    // Определяем мобильное устройство
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 1024);
      // Автоматически открываем сайдбар на мобильных при загрузке
      if (window.innerWidth < 1024 && !selectedFile) {
        setSidebarOpen(true);
      }
    };

    checkMobile();
    window.addEventListener('resize', checkMobile);

    // Обработка кнопки назад
    const handlePopState = (e: PopStateEvent) => {
      if (selectedFile && isMobile) {
        e.preventDefault();
        setSelectedFile(null);
        setFileContent('');
        setSidebarOpen(true);
        window.history.pushState(null, '', window.location.pathname);
      }
    };

    window.addEventListener('popstate', handlePopState);

    return () => {
      window.removeEventListener('resize', checkMobile);
      window.removeEventListener('popstate', handlePopState);
    };
  }, [selectedFile, isMobile]);

  const fetchDocFiles = async () => {
    try {
      // Импортируем AuthService для получения заголовков с токеном
      const { AuthService } = await import('@/services/auth');
      const headers = await AuthService.getAuthHeaders();
      
      const response = await fetch('/api/v1/docs/files', {
        headers,
        credentials: 'include'
      });
      
      if (response.ok) {
        const data: DocFilesResponse = await response.json();
        setFiles(data.data?.files || []);
      } else {
        console.error('Failed to fetch doc files:', response.status, response.statusText);
      }
    } catch (error) {
      console.error('Failed to fetch doc files:', error);
    }
  };

  const fetchFileContent = async (path: string) => {
    setLoading(true);
    try {
      // Импортируем AuthService для получения заголовков с токеном
      const { AuthService } = await import('@/services/auth');
      const headers = await AuthService.getAuthHeaders();
      
      const response = await fetch(
        `/api/v1/docs/content?path=${encodeURIComponent(path)}`,
        {
          headers,
          credentials: 'include'
        }
      );
      
      if (response.ok) {
        const data: DocContentResponse = await response.json();
        setFileContent(data.data?.content || '');
        setSelectedFile(path);
        // Закрываем сайдбар на мобильных устройствах после выбора файла
        setSidebarOpen(false);
        // Добавляем состояние в историю для обработки кнопки назад
        if (isMobile) {
          window.history.pushState(
            { file: path },
            '',
            window.location.pathname
          );
        }
      } else {
        console.error('Failed to fetch file content:', response.status, response.statusText);
      }
    } catch (error) {
      console.error('Failed to fetch file content:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleDirectory = (path: string) => {
    setExpandedDirs((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(path)) {
        newSet.delete(path);
      } else {
        newSet.add(path);
      }
      return newSet;
    });
  };

  const renderFileTree = (items: DocFile[], level = 0) => {
    return items.map((item) => {
      const isExpanded = expandedDirs.has(item.path || '');

      if (item.type === 'directory') {
        return (
          <div key={item.path}>
            <div
              className={`flex items-center gap-2 px-2 py-1 hover:bg-base-200 rounded cursor-pointer`}
              style={{ paddingLeft: `${level * 16 + 8}px` }}
              onClick={() => toggleDirectory(item.path || '')}
            >
              <ChevronRightIcon
                className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-90' : ''}`}
              />
              {isExpanded ? (
                <FolderOpenIcon className="w-5 h-5" />
              ) : (
                <FolderIcon className="w-5 h-5" />
              )}
              <span className="text-sm font-medium">{item.name}</span>
            </div>
            {isExpanded && item.children && (
              <div>{renderFileTree(item.children, level + 1)}</div>
            )}
          </div>
        );
      }

      return (
        <div
          key={item.path}
          className={`flex items-center gap-2 px-2 py-1 hover:bg-base-200 rounded cursor-pointer ${
            selectedFile === item.path ? 'bg-base-200' : ''
          }`}
          style={{ paddingLeft: `${level * 16 + 32}px` }}
          onClick={() => fetchFileContent(item.path || '')}
        >
          <DocumentTextIcon className="w-5 h-5" />
          <span className="text-sm">{item.name}</span>
        </div>
      );
    });
  };

  return (
    <AdminGuard>
      <div className="min-h-screen bg-base-200">
        {/* Mobile header */}
        <div className="lg:hidden flex items-center justify-between p-4 bg-base-100 shadow-lg">
          <div className="flex items-center gap-2">
            {selectedFile && (
              <button
                onClick={() => {
                  setSelectedFile(null);
                  setFileContent('');
                  setSidebarOpen(true);
                }}
                className="btn btn-ghost btn-sm btn-circle"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </button>
            )}
            <h1 className="text-xl font-bold">
              {selectedFile ? selectedFile.split('/').pop() : 'Documentation'}
            </h1>
          </div>
          {!selectedFile && (
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="btn btn-ghost btn-sm"
            >
              {sidebarOpen ? (
                <XMarkIcon className="w-6 h-6" />
              ) : (
                <Bars3Icon className="w-6 h-6" />
              )}
            </button>
          )}
        </div>

        {/* Desktop header */}
        <div className="hidden lg:block container mx-auto px-4 py-8">
          <h1 className="text-3xl font-bold mb-6">
            {selectedFile
              ? `Documentation: ${selectedFile.split('/').pop()}`
              : 'Documentation'}
          </h1>
        </div>

        <div className="lg:container lg:mx-auto lg:px-4">
          <div className="flex flex-col lg:flex-row gap-0 lg:gap-6">
            {/* Sidebar with file tree */}
            <div
              className={`${
                sidebarOpen || (!selectedFile && isMobile) ? 'block' : 'hidden'
              } lg:block fixed inset-0 z-50 lg:relative lg:inset-auto lg:z-auto bg-base-100 lg:bg-transparent`}
            >
              <div className="h-screen lg:h-auto lg:w-80 lg:bg-base-100 lg:rounded-lg lg:shadow-lg p-4 lg:h-[calc(100vh-12rem)] overflow-y-auto">
                {/* Mobile close button */}
                <div className="flex items-center justify-between mb-4 lg:hidden">
                  <h2 className="text-lg font-semibold">Files</h2>
                  <button
                    onClick={() => setSidebarOpen(false)}
                    className="btn btn-ghost btn-sm"
                  >
                    <XMarkIcon className="w-6 h-6" />
                  </button>
                </div>
                {/* Desktop title */}
                <h2 className="hidden lg:block text-lg font-semibold mb-4">
                  Files
                </h2>
                {files.length > 0 ? (
                  <div className="space-y-1">{renderFileTree(files)}</div>
                ) : (
                  <p className="text-sm text-base-content/60">
                    No documentation files found
                  </p>
                )}
              </div>
            </div>

            {/* Content area */}
            <div className="flex-1 bg-base-100 rounded-lg shadow-lg p-4 lg:p-6 min-h-[calc(100vh-4rem)] lg:h-[calc(100vh-12rem)] overflow-y-auto">
              {loading ? (
                <div className="flex items-center justify-center h-full">
                  <span className="loading loading-spinner loading-lg"></span>
                </div>
              ) : selectedFile ? (
                <div className="prose prose-sm md:prose-base lg:prose-lg max-w-none">
                  <ReactMarkdown
                    remarkPlugins={[remarkGfm]}
                    components={{
                      code({ className, children, ...props }: any) {
                        const match = /language-(\w+)/.exec(className || '');
                        const inline = !children?.toString().includes('\n');
                        return !inline && match ? (
                          <div className="bg-base-200 rounded-lg p-4 my-4 overflow-x-auto border border-base-300">
                            <SyntaxHighlighter
                              style={oneLight}
                              language={match[1]}
                              PreTag="div"
                              customStyle={{
                                backgroundColor: 'transparent',
                                padding: '0',
                                margin: '0',
                                overflow: 'visible',
                              }}
                              {...props}
                            >
                              {String(children).replace(/\n$/, '')}
                            </SyntaxHighlighter>
                          </div>
                        ) : (
                          <code
                            className="bg-base-200 text-base-content px-1 py-0.5 rounded text-sm"
                            {...props}
                          >
                            {children}
                          </code>
                        );
                      },
                      table({ children }) {
                        return (
                          <div className="overflow-x-auto">
                            <table className="table table-zebra">
                              {children}
                            </table>
                          </div>
                        );
                      },
                    }}
                  >
                    {fileContent}
                  </ReactMarkdown>
                </div>
              ) : (
                <div className="flex items-center justify-center h-full text-base-content/60">
                  <p className="text-center px-4">
                    {sidebarOpen || (!selectedFile && isMobile)
                      ? 'Select a file from the list'
                      : 'Click the menu button to browse documentation files'}
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </AdminGuard>
  );
}
