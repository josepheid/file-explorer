import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Login } from './Login';
import { RouterProvider, createMemoryRouter } from 'react-router';

describe('Login', () => {
  beforeEach(() => {
    vi.spyOn(window, 'fetch');
  });

  const renderWithRouter = () => {
    const router = createMemoryRouter([
      {
        path: '/',
        element: <Login />,
      },
    ]);
    return render(<RouterProvider router={router} />);
  };

  it('renders login form', () => {
    renderWithRouter();
    expect(screen.getByLabelText(/username/i)).toBeDefined();
    expect(screen.getByLabelText(/password/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /login/i })).toBeDefined();
  });

  it('shows error on failed login', async () => {
    vi.mocked(fetch).mockRejectedValueOnce(new Error('Failed to fetch'));
    renderWithRouter();

    fireEvent.change(screen.getByLabelText(/username/i), {
      target: { value: 'testuser' },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' },
    });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(screen.getByText(/Failed to fetch/i)).toBeDefined();
    });
  });

  it('disables button during login attempt', async () => {
    vi.mocked(fetch).mockImplementationOnce(
      () =>
        new Promise(resolve => setTimeout(() => resolve(new Response()), 100))
    );
    renderWithRouter();

    const button = screen.getByRole('button', { name: /login/i });
    fireEvent.change(screen.getByLabelText(/username/i), {
      target: { value: 'testuser' },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' },
    });

    fireEvent.click(button);
    expect(button).toBeDisabled();
    expect(button).toHaveTextContent(/loading/i);

    await waitFor(() => {
      expect(button).not.toBeDisabled();
    });
  });

  it('navigates on successful login', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(new Response());
    renderWithRouter();

    fireEvent.change(screen.getByLabelText(/username/i), {
      target: { value: 'testuser' },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' },
    });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(window.location.pathname).toBe('/');
    });
  });
});
