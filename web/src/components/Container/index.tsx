import styled from 'styled-components';

export const Container = styled.div<{ $minHeight?: string }>`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0rem 2rem;
  min-height: ${props => props.$minHeight || '100vh'};
  background-color: var(--guinness-black);
  color: var(--guinness-cream);
`;
