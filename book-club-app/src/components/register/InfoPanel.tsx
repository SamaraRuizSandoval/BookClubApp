import './infoPanel.css';
export function InfoPanel() {
  return (
    <>
      <aside className="inner-left-panel" aria-label="Why join BookClub">
        <div className="panel-brand">
          <h2 className="panel-tagline">
            Reading is better
            <br />
            with <em>the right people.</em>
          </h2>
          <p className="panel-desc">
            Join thousands of readers who track, review, and discuss the books
            they love — all in one cozy corner of the internet.
          </p>
        </div>

        <div className="mini-books" aria-hidden="true">
          <svg
            width="120"
            height="160"
            viewBox="0 0 120 160"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <rect
              x="8"
              y="10"
              width="72"
              height="100"
              rx="3"
              fill="#1a3a5c"
              transform="rotate(-8 8 10)"
            />
            <rect
              x="8"
              y="10"
              width="10"
              height="100"
              rx="3"
              fill="#0f2540"
              transform="rotate(-8 8 10)"
            />
            <rect
              x="40"
              y="20"
              width="68"
              height="96"
              rx="3"
              fill="#2c1a0e"
              transform="rotate(5 40 20)"
            />
            <rect
              x="40"
              y="20"
              width="10"
              height="96"
              rx="3"
              fill="#1a0e07"
              transform="rotate(5 40 20)"
            />
            <rect
              x="22"
              y="30"
              width="70"
              height="98"
              rx="3"
              fill="#1a2e1a"
              transform="rotate(-2 22 30)"
            />
            <rect
              x="22"
              y="30"
              width="10"
              height="98"
              rx="3"
              fill="#0e1e0e"
              transform="rotate(-2 22 30)"
            />
          </svg>
        </div>

        <ul className="panel-features" aria-label="Features">
          <li>
            <span className="feat-icon" aria-hidden="true">
              {' '}
              📖{' '}
            </span>
            <span>
              <strong>Track every book</strong> you've read, are reading, or
              want to read next.
            </span>
          </li>
          <li>
            <span className="feat-icon" aria-hidden="true">
              {' '}
              💬{' '}
            </span>
            <span>
              <strong>Write reviews</strong> and spark real conversations with
              your circle.
            </span>
          </li>
          <li>
            <span className="feat-icon" aria-hidden="true">
              {' '}
              ✨{' '}
            </span>
            <span>
              <strong>Discover new reads</strong> curated by people who get your
              taste.
            </span>
          </li>
        </ul>

        <p className="panel-footer" aria-label="Free, no credit card needed">
          Free forever. No credit card required.
          <br />
          <a href="/privacy">Privacy Policy</a> ·{' '}
          <a href="/terms">Terms of Service</a>
        </p>
      </aside>
    </>
  );
}
