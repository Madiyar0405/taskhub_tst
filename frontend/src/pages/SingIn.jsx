import React, {useState, useContext} from 'react';
import { useNavigate, Link } from 'react-router-dom';
import AuthContext from '../contexts/AuthContext'; // Import AuthContext
import Navbar from '../components/navbar/Navbar.jsx';
import Footer from '../components/footer/Footer.jsx';
import style from "../styles/LoginSignup.module.css";

const SignIn = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const { logIn, user } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const result = await logIn(username, password);

        if (result.success) {
            navigate('/');
        } else {
            setError(result.message);
        }
    };



    return (
        <>
            <Navbar />
            <div className={style.login_container}>
                <h1>Were happy your back!</h1>
                <div className={style.login_body}>
                    <form onSubmit={handleSubmit}>
                        {error && <div className={style.error_message}>{error}</div>}
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Username"
                            required
                        />
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Password"
                            required
                        />
                        <button type="submit">Login</button>
                        <div className={style.login_footer}>
                            <Link to="/signup" className={style.signup}>Sign Up</Link>
                            <Link to="/forgot-password" className={style.forgot_password}>Forgot your password?</Link>
                        </div>
                    </form>
                </div>
            </div>
            <Footer />
        </>
    );
};

export default SignIn;
