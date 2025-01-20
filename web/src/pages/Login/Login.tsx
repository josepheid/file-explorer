import React from 'react';
import { Title, Button, Center, Container } from '../../components';
import styled from 'styled-components';
import { useNavigate } from 'react-router';

const Form = styled.form`
  background: var(--guinness-cream);
  padding: 2rem;
  border-radius: 0.5rem;
  box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.1);
  width: 90%;
  max-width: 20rem;
`;

const FormGroup = styled.div`
  margin-bottom: 1rem;
`;

const Label = styled.label`
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-size: 1rem;
`;

const Input = styled.input`
  width: 100%;
  padding: 0.75rem;
  border: 0.0625rem solid #ddd;
  border-radius: 0.25rem;
  font-size: 1rem;
  box-sizing: border-box;
`;

const ErrorMessage = styled.div`
  color: #dc3545;
  margin-bottom: 1rem;
  text-align: center;
  font-size: 0.875rem;
`;

export function Login() {
  const [error, setError] = React.useState<string>('');
  const [isLoading, setIsLoading] = React.useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const formData = new FormData(form);

    setError('');
    setIsLoading(true);

    try {
      const response = await fetch('/api/v1/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: formData.get('username'),
          password: formData.get('password'),
        }),
      });

      if (response.status === 401) {
        throw new Error('Invalid credentials');
      }
      if (!response.ok) {
        throw new Error('Login failed. Please try again later.');
      }

      await navigate('/browse');
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container>
      <Form onSubmit={e => void handleSubmit(e)}>
        <Title>Please login</Title>
        {error && <ErrorMessage>{error}</ErrorMessage>}
        <FormGroup>
          <Label htmlFor="username">Username</Label>
          <Input type="text" id="username" name="username" required />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="password">Password</Label>
          <Input type="password" id="password" name="password" required />
        </FormGroup>
        <Center>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Login'}
          </Button>
        </Center>
      </Form>
    </Container>
  );
}
