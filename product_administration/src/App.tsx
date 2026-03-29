import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import CategoryManagement from './views/CategoryManagement';
import SubCategoryManagement from './views/SubCategoryManagement';
import Dashboard from './views/Dashboard';
import RequireAuth from './components/RequireAuth';
import ProductTypeManagement from './views/ProductTypeManagement';
import StoringConditionManagement from './views/StoringConditionManagement';
import BrandManagement from './views/BrandManagement';
import ProductManagement from './views/ProductManagement';
import SearchProductManagement from './views/SearchProductManagement';
import { AuthProvider } from './context/AuthContext';

const Login = () => <h1>Login Page</h1>;

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route element={<RequireAuth />}>
            <Route path="/" element={<Layout />}>
              <Route index element={<Dashboard />} />
              <Route path="categories" element={<CategoryManagement />} />
              <Route path="subcategories" element={<SubCategoryManagement/>} />
              <Route path="storing-condition" element={<StoringConditionManagement />} />
              <Route path="product-types" element={<ProductTypeManagement />} />
              <Route path="products" element={<ProductManagement/>} />
              <Route path="brands" element={<BrandManagement />} />
              <Route path="search" element={<SearchProductManagement/>} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}