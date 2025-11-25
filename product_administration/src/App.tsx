import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/layout/Layout';
import CategoryManagement from './views/CategoryManagement';
import SubCategoryManagement from './views/SubCategoryManagement';
import Dashboard from './views/Dashboard';
import { ProductPlaceholder, ProductTypePlaceholder, BrandPlaceholder, SearchPlaceholder } from './views/Placeholders';

import RequireAuth from './components/RequireAuth'; 

const Login = () => <h1>Login Page</h1>;


export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        
        {/* Public Routes - Accessible without login */}
        <Route path="/login" element={<Login />} />


        {/* Protected Routes 
          This parent route uses RequireAuth. Any nested route inside it will 
          first run the logic in RequireAuth (the authentication check).
        */}
        <Route element={<RequireAuth />}>
          <Route path="/" element={<Layout />}>
            {/* The routes below are protected: */}
            <Route index element={<Dashboard />} />
            <Route path="categories" element={<CategoryManagement />} />
            <Route path="subcategories" element={<SubCategoryManagement/>} />
            <Route path="products" element={<ProductPlaceholder />} />
            <Route path="product-types" element={<ProductTypePlaceholder />} />
            <Route path="brands" element={<BrandPlaceholder />} />
            <Route path="search" element={<SearchPlaceholder />} />
            
            {/* Catch-all route for any undefined path */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Route>

      </Routes>
    </BrowserRouter>
  );
}