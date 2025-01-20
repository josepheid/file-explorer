import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { Browse } from './Browse';
import { BrowserRouter } from 'react-router';

// Mock directory response
const mockDirectoryContents = {
  name: 'test',
  type: 'dir',
  size: 0,
  contents: [
    { name: 'folder1', type: 'dir', size: 0 },
    { name: 'folder2', type: 'dir', size: 0 },
    { name: 'test.txt', type: 'file', size: 1024 },
    { name: 'example.pdf', type: 'file', size: 2048 },
  ],
};

const renderBrowse = () => {
  return render(
    <BrowserRouter>
      <Browse />
    </BrowserRouter>
  );
};

describe('Browse', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn(() =>
      Promise.resolve(
        new Response(JSON.stringify(mockDirectoryContents), {
          headers: { 'Content-Type': 'application/json' },
        })
      )
    );
    // Reset URL to root
    window.history.pushState({}, '', '/');
  });

  it('renders the browse page', async () => {
    renderBrowse();
    expect(screen.getByText('Browse Files')).toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByText('folder1')).toBeInTheDocument();
    });
  });

  it('displays directory contents correctly', async () => {
    renderBrowse();

    await waitFor(() =>
      expect(screen.getByText('folder1')).toBeInTheDocument()
    );
    expect(screen.getByText('folder2')).toBeInTheDocument();
    expect(screen.getByText('test.txt')).toBeInTheDocument();
    expect(screen.getByText('example.pdf')).toBeInTheDocument();
  });

  it('sorts directories first by default', async () => {
    renderBrowse();

    await waitFor(() => {
      const fileItems = screen.getAllByTestId('file-item-name');
      const fileNames = fileItems.map(item => item.textContent);
      expect(fileNames.indexOf('folder1')).toBeLessThan(
        fileNames.indexOf('test.txt')
      );
      expect(fileNames.indexOf('folder2')).toBeLessThan(
        fileNames.indexOf('example.pdf')
      );
    });
  });

  it('filters contents based on search term', async () => {
    renderBrowse();
    const searchInput = screen.getByPlaceholderText(
      'Search files and folders...'
    );

    fireEvent.change(searchInput, { target: { value: 'folder' } });

    await waitFor(() => {
      expect(screen.getByText('folder1')).toBeInTheDocument();
      expect(screen.getByText('folder2')).toBeInTheDocument();
      expect(screen.queryByText('test.txt')).not.toBeInTheDocument();
    });
  });

  it('handles sorting by different fields', async () => {
    renderBrowse();
    const sortSelect = await screen.findByRole('combobox');

    fireEvent.change(sortSelect, { target: { value: 'name-desc' } });

    await waitFor(() => {
      const fileItems = screen.getAllByTestId('file-item-name');
      const fileNames = fileItems.map(item => item.textContent);
      expect(fileNames.indexOf('test.txt')).toBeLessThan(
        fileNames.indexOf('example.pdf')
      );
    });
  });

  it('shows error message when fetch fails', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        status: 500,
      })
    );

    renderBrowse();
    await waitFor(() => {
      expect(
        screen.getByText('Failed to load directory contents')
      ).toBeInTheDocument();
    });
  });

  it('updates URL when sorting or filtering', async () => {
    renderBrowse();
    const searchInput = await screen.findByPlaceholderText(
      'Search files and folders...'
    );

    fireEvent.change(searchInput, { target: { value: 'test' } });

    expect(window.location.search).toContain('filter=test');
  });

  it('handles breadcrumb navigation', async () => {
    renderBrowse();
    await waitFor(() => {
      expect(screen.getByText('root')).toBeInTheDocument();
    });

    const folder = await screen.findByText('folder1');
    fireEvent.click(folder);

    await waitFor(() => {
      expect(window.location.pathname).toContain('/folder1');
    });
  });

  it('formats file sizes correctly', async () => {
    renderBrowse();
    await waitFor(() => {
      expect(screen.getByText('1.0 KB')).toBeInTheDocument(); // 1024 bytes
      expect(screen.getByText('2.0 KB')).toBeInTheDocument(); // 2048 bytes
    });
  });

  it('handles unauthorized access', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        status: 401,
      })
    );

    renderBrowse();
    await waitFor(() => {
      expect(window.location.pathname).toBe('/login');
    });
  });
});
