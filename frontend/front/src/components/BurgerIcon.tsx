type BurgerIconProps = {
  onClick: () => void
}


const BurgerIcon = ({onClick}: BurgerIconProps) => {
  return (
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        aria-hidden
        id="burger-icon"
        onClick={() => onClick()}
      >
          <path
            d="M4 7H20"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <path
            d="M4 12H20"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <path
            d="M4 17H20"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
          />
      </svg>
)
}

export default BurgerIcon
