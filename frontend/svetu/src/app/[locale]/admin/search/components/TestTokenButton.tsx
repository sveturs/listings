'use client';

export default function TestTokenButton() {
  const addTestToken = () => {
    // Добавляем тестовый JWT токен для админа
    const testToken =
      'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzU1ODY0NzY4LCJpYXQiOjE3NTU3NzgzNjgsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoxfQ.y1JB88pvgib0e5QsS8sZYsSmAHX0fccMPvxCvzPh2x4';
    localStorage.setItem('token', testToken);
    window.location.reload();
  };

  return (
    <button onClick={addTestToken} className="btn btn-sm btn-primary">
      Добавить тестовый токен
    </button>
  );
}
