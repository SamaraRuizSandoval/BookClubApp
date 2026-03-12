// import './stars.css';

export function StarsBackground() {
  const stars = Array.from({ length: 60 }).map((_, i) => {
    const size = Math.random() * 2.5 + 0.5;

    return (
      <div
        key={i}
        className="star"
        style={{
          width: `${size}px`,
          height: `${size}px`,
          top: `${Math.random() * 100}%`,
          left: `${Math.random() * 100}%`,
          animationDelay: `${(Math.random() * 5).toFixed(1)}s`,
          opacity: Math.random() * 0.4 + 0.05,
        }}
      />
    );
  });

  return <div className="stars-bg">{stars}</div>;
}
