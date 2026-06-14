import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import ProjectCreate from './pages/ProjectCreate';
import ProjectBuild from './pages/ProjectBuild';
import ProjectGuide from './pages/ProjectGuide';
import ContentManager from './pages/ContentManager';
import ActivityManager from './pages/ActivityManager';
import Billing from './pages/Billing';
import ErrorBoundary from './components/ErrorBoundary';

function App() {
  return (
    <ConfigProvider locale={zhCN} theme={{ token: { colorPrimary: '#4caf50' } }}>
      <BrowserRouter>
        <ErrorBoundary>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/projects/new" element={<ProjectCreate />} />
            <Route path="/projects/:id/build" element={<ProjectBuild />} />
            <Route path="/projects/:id/guide" element={<ProjectGuide />} />
            <Route path="/content" element={<ContentManager />} />
            <Route path="/activities" element={<ActivityManager />} />
            <Route path="/billing" element={<Billing />} />
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </ErrorBoundary>
      </BrowserRouter>
    </ConfigProvider>
  );
}

export default App;
