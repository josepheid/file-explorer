import styled from 'styled-components';

export const Subtitle = styled.h2<{
  color?: string;
  fontSize?: string;
  marginBottom?: string;
}>`
  font-size: ${props => props.fontSize || '1.25rem'};
  text-align: center;
  margin-bottom: ${props => props.marginBottom || '2rem'};
  color: ${props => props.color || 'var(--guinness-cream)'};
`;
