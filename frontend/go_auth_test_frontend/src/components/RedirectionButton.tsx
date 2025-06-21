import type { ReactNode } from "react"
import { useNavigate } from "react-router-dom"

type RedirectionButtonProps = {
  to: `/${string}`,
  variation: "light" | "dark",
  children: ReactNode
}
const RedirectionButton = (props: RedirectionButtonProps) => {
  const navigate = useNavigate()
  return (
    <button onClick={() => navigate(props.to)}>
      {props.children}
    </button>
  )
}

export default RedirectionButton
