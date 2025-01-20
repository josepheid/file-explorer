import styled from 'styled-components';

export const Title = styled.h1<{
  color?: string;
  fontSize?: string;
  marginBottom?: string;
}>`
  font-size: ${props => props.fontSize || '1.5rem'};
  text-align: center;
  margin-bottom: ${props => props.marginBottom || '1.5rem'};
  color: ${props => props.color || 'var(--guinness-gold)'};
`;
