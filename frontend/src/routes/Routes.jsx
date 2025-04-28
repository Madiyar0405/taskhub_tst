import React, { useContext } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import Home from "../pages/Home";
import SignIn from "../pages/SingIn";
import SignUp from "../pages/SignUp";
import Verification from "../pages/Verification";
import Profile from "../pages/Profile";
import AuthContext from "../contexts/AuthContext";
import GitHubCallback from "../components/GitHubCallback";
import CompaniesPage from "../pages/CompaniesPage";
import CompanyInfo from "../pages/CompanyPage";
import CompanyIssue from "../pages/Issue";

export function AppRoutes() {
    const { user, isLoading } = useContext(AuthContext);

    if (isLoading) return <div>Loading...</div>; // ✅ Показываем "Загрузка..." вместо редиректа

    return (
        <Routes>
            {/* Public routes */}
            <Route path="/" element={<Home />} />
            <Route path="/signin" element={<SignIn />} />
            <Route path="/signup" element={<SignUp />} />
            <Route path="/signup/verification" element={<Verification />} />
            <Route path="/github/callback" element={<GitHubCallback />} />

            {/* Private routes */}
            <Route path="/profile/:developerID" element={user ? <Profile /> : <Navigate to="/signin" replace />} />
            <Route path="/companies" element={user ? <CompaniesPage /> : <Navigate to="/signin" replace />} />
            <Route path="/companies/:companyId" element={user ? <CompanyInfo /> : <Navigate to="/signin" replace />} />
            <Route path="/issues/:issueId" element={<CompanyIssue />} />

            {/* Fallback route */}
            <Route path="*" element={<h1>Page Not Found</h1>} />
        </Routes>
    );
}
