import { NavLink } from "react-router-dom"

const VerifyEmail = () => {
  return (
    <div>
      <h2>Your account was created</h2>
      <p>Click the verification link that was sent to your email address</p>
      <NavLink to="/login">Go To Login Page</NavLink>
    </div>
  )
}

export default VerifyEmail
