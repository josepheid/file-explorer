import { useEffect, useState } from 'react';
import { useNavigate, useLocation, useSearchParams } from 'react-router';
import styled from 'styled-components';
import { Title, Container, Navbar } from '../../components';
import { Link } from 'react-router';

interface FileInfo {
  name: string;
  type: string;
  size: number;
}

interface DirectoryContents {
  name: string;
  type: string;
  size: number;
  contents: FileInfo[];
}

type SortField = 'name' | 'type' | 'size';
type SortDirection = 'asc' | 'desc';

const FileExplorer = styled.div`
  padding: 1rem;
  background-color: var(--guinness-cream);
  border: 1px solid black;
  border-radius: 4px;
  width: 95%;

  @media (min-width: 768px) {
    width: 80%;
    padding: 2rem;
  }
`;
const FileList = styled.div`
  border: 1px solid black;
  background-color: white;
`;

const FileListHeader = styled.div`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #e0e0e0;
  background-color: var(--guinness-gold);
  color: black;
  font-weight: bold;
`;

const FileItem = styled.div<{ $isDirectory: boolean }>`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #e0e0e0;
  color: black;
  cursor: ${props => (props.$isDirectory ? 'pointer' : 'default')};

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background-color: #f0f0f0;
  }
`;

const FileName = styled.span`
  flex: 1;
`;

const FileType = styled.span`
  margin: 0 1rem;
`;

const FileSize = styled.span`
  min-width: 100px;
  text-align: right;
`;

const Breadcrumbs = styled.div`
  display: flex;
  align-items: center;
  margin: 0 0 1rem 0;
  font-size: 1.5rem;
`;

const BreadcrumbLink = styled(Link)`
  color: var(--guinness-gold);
  text-decoration: none;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
`;

const BreadcrumbItem = styled.span`
  color: var(--guinness-gold);
  text-decoration: none;
`;

const BreadcrumbSeparator = styled.span`
  margin: 0 0.5rem;
  color: var(--guinness-gold);
`;

const ControlBar = styled.div`
  margin: 1rem 0;
  @media (min-width: 768px) {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
`;

const Select = styled.select`
  padding: 0.5rem;
  border-radius: 4px;
  border: 1px solid black;
  background-color: white;
  color: black;
  cursor: pointer;

  &:hover {
    border-color: black;
  }
`;

const SearchBar = styled.input`
  padding: 0.5rem;
  border-radius: 4px;
  border: 1px solid black;
  background-color: white;
  color: black;
  margin-right: auto; // This pushes the sort dropdown to the right

  &:focus {
    outline: none;
    border-color: black;
  }
`;

export function Browse() {
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams, setSearchParams] = useSearchParams();
  const [directory, setDirectory] = useState<DirectoryContents | null>(null);
  const [error, setError] = useState<string>('');

  // Get path from URL, defaulting to '/'
  const currentPath = decodeURIComponent(
    location.pathname.replace('/browse', '') || '/'
  );

  // Initialize state from URL params
  const [sortField, setSortField] = useState<SortField>(
    (searchParams.get('sort') as SortField) || 'type'
  );
  const [sortDirection, setSortDirection] = useState<SortDirection>(
    (searchParams.get('direction') as SortDirection) || 'asc'
  );
  const [searchTerm, setSearchTerm] = useState(
    searchParams.get('filter') || ''
  );

  // Update URL when sort or filter changes
  const updateUrlParams = (
    newSort?: SortField,
    newDirection?: SortDirection,
    newFilter?: string
  ) => {
    const params = new URLSearchParams(searchParams);

    if (newSort) params.set('sort', newSort);
    if (newDirection) params.set('direction', newDirection);
    if (newFilter !== undefined) {
      if (newFilter) {
        params.set('filter', newFilter);
      } else {
        params.delete('filter');
      }
    }

    setSearchParams(params);
  };

  const handleSortChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const [field, direction] = event.target.value.split('-') as [
      SortField,
      SortDirection,
    ];
    setSortField(field);
    setSortDirection(direction);
    updateUrlParams(field, direction);
  };

  const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    setSearchTerm(value);
    updateUrlParams(undefined, undefined, value);
  };

  useEffect(() => {
    const fetchDirectory = async () => {
      try {
        const response = await fetch(
          `/api/v1/browse?path=${encodeURIComponent(currentPath)}`,
          {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
          }
        );

        if (response.status === 401) {
          await navigate('/login');
          return;
        }

        if (response.status === 404) {
          setError('Directory not found');
          return;
        }

        if (!response.ok) {
          throw new Error('Failed to fetch directory contents');
        }

        const data = (await response.json()) as DirectoryContents;
        setDirectory(data);
        setError('');
      } catch (err) {
        setError('Failed to load directory contents');
        console.error(err);
      }
    };

    void fetchDirectory();
  }, [currentPath, navigate]);

  const handleFileClick = async (file: FileInfo) => {
    if (file.type === 'dir') {
      const newPath =
        currentPath === '/' ? `/${file.name}` : `${currentPath}/${file.name}`;
      await navigate(`/browse${newPath}`);
    }
  };

  const formatFileSize = (size: number): string => {
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
    return `${(size / (1024 * 1024)).toFixed(1)} MB`;
  };

  const renderBreadcrumbs = () => {
    const pathParts = currentPath.split('/').filter(Boolean);
    return (
      <Breadcrumbs>
        {pathParts.length === 0 ? (
          <BreadcrumbItem>root</BreadcrumbItem>
        ) : (
          <BreadcrumbLink to="/browse/">root</BreadcrumbLink>
        )}
        {pathParts.map((part, index) => {
          const path = `/browse/${pathParts.slice(0, index + 1).join('/')}`;
          return (
            <span key={path}>
              <BreadcrumbSeparator>/</BreadcrumbSeparator>
              {index === pathParts.length - 1 ? (
                <BreadcrumbItem>{part}</BreadcrumbItem>
              ) : (
                <BreadcrumbLink to={path}>{part}</BreadcrumbLink>
              )}
            </span>
          );
        })}
      </Breadcrumbs>
    );
  };

  const filteredAndSortedContents = directory?.contents
    .filter(file => file.name.toLowerCase().includes(searchTerm.toLowerCase()))
    .sort((a, b) => {
      let comparison = 0;

      // Then apply the selected sort
      switch (sortField) {
        case 'name':
          comparison = a.name.localeCompare(b.name);
          break;
        case 'type':
          comparison = a.type.localeCompare(b.type);
          break;
        case 'size':
          comparison = a.size - b.size;
          break;
      }

      return sortDirection === 'asc' ? comparison : -comparison;
    });

  return (
    <>
      <Navbar />
      <Container $minHeight="calc(100vh - 5.0625rem)">
        <Title>Browse Files</Title>

        <FileExplorer>
          {renderBreadcrumbs()}

          <ControlBar>
            <SearchBar
              type="text"
              placeholder="Search files and folders..."
              value={searchTerm}
              onChange={handleSearch}
            />
            <Select
              value={`${sortField}-${sortDirection}`}
              onChange={handleSortChange}
            >
              <option value="type-asc">Type (Directories first)</option>
              <option value="type-desc">Type (Files first)</option>
              <option value="name-asc">Name (A to Z)</option>
              <option value="name-desc">Name (Z to A)</option>
              <option value="size-asc">Size (Smallest first)</option>
              <option value="size-desc">Size (Largest first)</option>
            </Select>
          </ControlBar>

          {error && (
            <div style={{ color: 'red', margin: '1rem 0' }}>{error}</div>
          )}

          {directory && (
            <FileList>
              <FileListHeader>
                <FileName>Name</FileName>
                <FileType>Type</FileType>
                <FileSize>Size</FileSize>
              </FileListHeader>
              {filteredAndSortedContents?.map(file => (
                <FileItem
                  key={file.name}
                  $isDirectory={file.type === 'dir'}
                  onClick={() => void handleFileClick(file)}
                  data-testid={'file-item'}
                >
                  <FileName data-testid={'file-item-name'}>
                    {file.name}
                  </FileName>
                  <FileType>
                    {file.type === 'dir' ? 'Directory' : 'File'}
                  </FileType>
                  <FileSize>{formatFileSize(file.size)}</FileSize>
                </FileItem>
              ))}
            </FileList>
          )}
        </FileExplorer>
      </Container>
    </>
  );
}
