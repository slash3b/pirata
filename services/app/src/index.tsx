import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import About from './About';
// import Next from './Next';

import reportWebVitals from './reportWebVitals';
import { BrowserRouter, Route, Routes } from 'react-router';

const root = document.getElementById("root") as HTMLElement;

ReactDOM.createRoot(root).render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<App />} />
      <Route path="/about" element={<About />} >
        {/* <Route path="/next" element={<Next />} /> */}
      </Route>
    </Routes>
  </BrowserRouter>

)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
