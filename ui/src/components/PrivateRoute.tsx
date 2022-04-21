import { Navigate, useLocation } from "react-router-dom";
import AuthService from "../services/auth-service";

const PrivateElement = ({ children }: any) => {
  let location = useLocation();

  return isAuthorized() ? (
    children
  ) : (
    <Navigate to="/login" state={{ from: location }} />
  );
};

function isAuthorized() {
  const token = AuthService.getToken();

  try {
    if (token) {
      const decodedJwt = JSON.parse(atob(token.split(".")[1]));
      if (decodedJwt.exp * 1000 < Date.now()) {
        AuthService.logout();
        return false;
      }
    } else {
      return false;
    }

    return true;
  } catch (e) {
    console.error("isAuthorized Error:", e);
    return false;
  }
}

export default PrivateElement;
