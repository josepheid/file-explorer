import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { Navbar } from './index';
import { MemoryRouter } from 'react-router';

describe('Navbar', () => {
  const renderWithRouter = (component: React.ReactNode) => {
    return render(<MemoryRouter>{component}</MemoryRouter>);
  };

  it('renders logo', () => {
    renderWithRouter(<Navbar />);
    expect(screen.getByText('Stout Explorer')).toBeDefined();
  });

  it('renders logout button', () => {
    renderWithRouter(<Navbar />);
    expect(screen.getByRole('button', { name: /logout/i })).toBeDefined();
  });
});
