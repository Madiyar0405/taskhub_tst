import React, { useState, useContext, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";
import AuthContext from "../contexts/AuthContext";
import style from "../styles/LoginSignup.module.css";

const SignUp = () => {
    const [email, setEmail] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [repeatPassword, setRepeatPassword] = useState('');
    const [error, setError] = useState('');
    const { signUp, verificationCodeSent } = useContext(AuthContext);
    const navigate = useNavigate();

    useEffect(() => {
        if (verificationCodeSent) {
            navigate('/signup/verification', { state: { email } });
        }
    }, [verificationCodeSent, email, navigate]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (password !== repeatPassword) {
            setError("Passwords do not match.");
            return;
        }
        if (password.length < 6) {
            setError("Password must be at least 6 characters long.");
            return;
        }
        setError('');
        const result = await signUp(email, password, username);

        console.log(result+ " result");
        if (result.success) {
            navigate('/login');
        } else {
            setError(result.message);
        }
    };



    return (
        <div>
            <Navbar />
            <div className={style.login_container}>
                <h1>Fill out the form</h1>
                <div className={style.login_body}>

                    <form onSubmit={handleSubmit}>
                        {error && <div className={style.error_message}>{error}</div>}

                        <input
                            type="text"
                            value={email}
                            name="email"
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="Email"
                            required
                        />
                        <input
                            type="text"
                            value={username}
                            name="username"
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Username"
                            required
                        />
                        <input
                            type="password"
                            value={password}
                            name="password"
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Password"
                            required
                        />
                        <input
                            type="password"
                            value={repeatPassword}
                            name="repeatPassword"
                            onChange={(e) => setRepeatPassword(e.target.value)}
                            placeholder="Repeat Password"
                            required
                        />
                        <button type="submit">Sign Up</button>
                        <div className={style.login_footer}>
                            <Link to="/signIn" className={style.sign_in}>Sign In</Link>
                            <Link to="/forgot-password" className={style.forgot_password}>Forgot your password?</Link>
                        </div>
                    </form>
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default SignUp;
