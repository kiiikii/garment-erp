//! unlock interactivity in browser
"use client"

import { useState } from "react";

export default function CustomerForm() {
  //! create react state hold the data as the user type
  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");

  //! function that running when user click Submit
  const handleSubmit = async (e: React.FormEvent) => {
    //! stop page to refreshing
    e.preventDefault();

    //! package into JSON
    const newCustomer = { name, phone, address };

    //! fire the POST request
    const res = await fetch("http://localhost:8080/customers", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(newCustomer),
    });

    if (res.ok) {
      alert("Customer added Successfully");

      //! clearing field
      setName("");
      setPhone("");
      setAddress("");

      //! forcing page to refresh to show new data table
      window.location.reload();
    } else {
      const errorText = await res.text();
      alert("Failed to add customer: " + errorText);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white p-6 rounded-lg shadow mb-8"
    >
      <h2 className="text-xl font-semibold mb-4 text-gray-800">
        Add new Customer
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <input
          type="text"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="border p-2 rounded text-black"
          required
        />
        <input
          type="text"
          placeholder="Phone"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className="border p-2 rounded text-black"
          required
        />
        <input
          type="text"
          placeholder="Address"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          className="border p-2 rounded text-black"
          required
        />
      </div>
      <button
        type="submit"
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors"
      >
        Save Customer
      </button>
    </form>
  );
}
