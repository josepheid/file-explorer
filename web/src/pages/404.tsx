import { Title, Container, StyledLink } from '../components';
export function NotFound() {
  return (
    <Container>
      <Title>Path not found</Title>
      <StyledLink to="/browse">Back to Home</StyledLink>
    </Container>
  );
}
