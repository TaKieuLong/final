import React, { useState } from 'react';
import { createProduct } from '../api';

const CreateProductForm = () => {
    const [formData, setFormData] = useState({
        name: '',
        description: '',
        price: '',
        image_url: '',
    });

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        await createProduct(formData);
        alert('Product created!');
    };

    return (
        <form onSubmit={handleSubmit}>
            <input name="name" placeholder="Name" onChange={handleChange} />
            <input name="description" placeholder="Description" onChange={handleChange} />
            <input name="price" placeholder="Price" onChange={handleChange} />
            <input name="image_url" placeholder="Image URL" onChange={handleChange} />
            <button type="submit">Create</button>
        </form>
    );
};

export default CreateProductForm;
