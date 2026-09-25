type BurgerIconProps = {
  onClick: () => void;
};

const BurgerIcon = ({onClick}: BurgerIconProps) => {
  return (
    <button
      type="button"
      id="burger-icon"
      className="btn-ghost"
      aria-label="Toggle menu"
      onClick={() => onClick()}
    >
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        aria-hidden
      >
          <path
            d="M4 7H20"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <path
            d="M4 12H20"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <path
            d="M4 17H20"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
          />
      </svg>
    </button>
)
}

export default BurgerIcon
