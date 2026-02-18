import React from 'react';

export const NotFound: React.FC = () => {
  return (
    <div style={{
      maxWidth: '1180px',
      minHeight: '100vh',
      margin: '0 auto',
      padding: '1rem',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      flexDirection: 'column',
      gap: '1rem'
    }}>
      <h1 style={{ fontSize: '4rem', color: '#E53E3E' }}>404</h1>
      <h2 style={{ fontSize: '2rem', color: '#2D3748' }}>Page Not Found</h2>
      <p style={{ color: '#718096' }}>The page you're looking for doesn't exist.</p>
      <a
        href="/"
        style={{
          padding: '0.5rem 1rem',
          backgroundColor: '#4299E1',
          color: 'white',
          border: 'none',
          borderRadius: '0.25rem',
          textDecoration: 'none'
        }}
      >
        Go Home
      </a>
    </div>
  );
};
