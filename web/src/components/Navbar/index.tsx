import styled from 'styled-components';
import { Button } from '../Button';
import { useState } from 'react';
import { useNavigate } from 'react-router';

const Nav = styled.nav`
  background-color: var(--guinness-cream);
  padding: 1rem 2rem;
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.2);
  height: 3rem;
`;

const NavContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 75rem;
  margin: 0 auto;
`;

const Logo = styled.span`
  color: var(--guinness-gold);
  font-size: 1.5rem;
  font-weight: bold;
`;

const NavLinks = styled.div`
  display: flex;
  gap: 2rem;
`;

const ErrorToast = styled.div`
  position: fixed;
  top: 1rem;
  right: 1rem;
  background-color: #dc3545;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 0.25rem;
  animation: fadeOut 3s forwards;

  @keyframes fadeOut {
    0% {
      opacity: 1;
    }
    70% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
`;

export function Navbar() {
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const logout = async () => {
    try {
      const response = await fetch('/api/v1/logout', { method: 'POST' });
      if (!response.ok) throw new Error('Logout failed');
      await navigate('/login');
    } catch (error) {
      setError('Failed to logout. Please try again.');
      setTimeout(() => setError(null), 3000);
      console.error('Error logging out', error);
    }
  };

  return (
    <>
      <Nav>
        <NavContent>
          <Logo>Stout Explorer</Logo>
          <NavLinks>
            <Button $secondary $isLink onClick={() => void logout()}>
              Logout
            </Button>
          </NavLinks>
        </NavContent>
      </Nav>
      {error && <ErrorToast>{error}</ErrorToast>}
    </>
  );
}
