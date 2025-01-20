import { Link } from 'react-router';
import styled from 'styled-components';

export const StyledLink = styled(Link)`
  color: var(--guinness-gold);
  text-decoration: none;
  cursor: pointer;
  &:hover {
    text-decoration: underline;
  }
`;
