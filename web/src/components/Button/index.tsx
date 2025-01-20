import styled from 'styled-components';

export const Button = styled.button<{
  $secondary?: boolean;
  $isLink?: boolean;
}>`
  width: auto;
  min-width: 8rem;
  padding: 0.75rem 1.5rem;
  background-color: ${props =>
    props.$secondary ? 'var(--guinness-cream)' : 'var(--guinness-black)'};
  color: ${props =>
    props.$secondary ? 'var(--guinness-gold)' : 'var(--guinness-cream)'};
  font-family: 'IM Fell DW Pica SC', serif;
  font-size: 1.25rem;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  transition: background-color 0.2s;

  ${props =>
    !props.$isLink &&
    `
    &:hover {
      background-color: var(--guinness-gold);
      color: var(--guinness-black);
    }
  `}

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;
