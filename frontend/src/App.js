import React from "react";
import ProductList from "./components/ProductList";
import CreateProductForm from "./components/CreateProductForm";

const App = () => {
  return (
    <div>
      <CreateProductForm />
      <ProductList />
    </div>
  );
};

export default App;
