"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function OrderForm() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [customers, setCustomers] = useState<{ id: number; name: string }[]>(
    [],
  );

  const [customerID, setCustomerID] = useState("");
  const [productionType, setProductionType] = useState("CMT");
  const [totalQuantity, setTotalQuantity] = useState("");
  const [files, setFiles] = useState<FileList | null>(null);

  const [sizes, setSizes] = useState<{ size: string; quantity: string }[]>([
    { size: "S", quantity: "" },
    { size: "M", quantity: "" },
    { size: "L", quantity: "" },
    { size: "XL", quantity: "" },
  ]);

  useEffect(() => {
    fetch("http://localhost:8080/customers")
      .then((res) => res.json())
      .then((data) => setCustomers(data))
      .catch((err) => console.error("Failed to load Customers", err));
  }, []);

  const handleSizeChange = (
    index: number,
    field: "size" | "quantity",
    value: string,
  ) => {
    const newSizes = [...sizes];
    newSizes[index][field] = value;
    setSizes(newSizes);

    if (field === "quantity") {
      const autoTotal = newSizes.reduce(
        (sum, current) => sum + (parseInt(current.quantity) || 0),
        0,
      );
      setTotalQuantity(autoTotal > 0 ? autoTotal.toString() : "");
    }
  };

  const addSizeRow = () => setSizes([...sizes, { size: "", quantity: "" }]);
  const removeSizeRow = (index: number) => {
    const newSizes = sizes.filter((_, i) => i !== index);
    setSizes(newSizes);
    const autoTotal = newSizes.reduce(
      (sum, current) => sum + (parseInt(current.quantity) || 0),
      0,
    );
    setTotalQuantity(autoTotal > 0 ? autoTotal.toString() : "");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!customerID || !totalQuantity)
      return alert("Please fill all the required fields");
    setLoading(true);

    try {
      //! base order
      const orderRes = await fetch("http://localhost:8080/orders", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          customer_id: parseInt(customerID),
          production_type: productionType,
          total_quantity: totalQuantity,
          status: "WAITING_FOR_SAMPLE",
          sizes: sizes
            .filter((s) => s.size.trim() !== "" && parseInt(s.quantity) > 0)
            .map((s) => ({
              size: s.size,
              quantity: parseInt(s.quantity),
            })),
        }),
      });

      if (!orderRes.ok) throw new Error("Failed to create order");
      const newOrder = await orderRes.json();

      //! upload images verified
      if (files && files.length > 0) {
        const formData = new FormData();
        formData.append("order_id", newOrder.id.toString());
        for (let i = 0; i < files.length; i++) {
          formData.append("tech_packs", files[i]);
        }
        const imageRes = await fetch("http://localhost:8080/orders/images", {
          method: "POST",
          body: formData,
        });

        if (!imageRes.ok) {
          const errText = await imageRes.text();
          throw new Error(
            "Order created, but images failed to upload: " + errText,
          );
        }
      }

      //! navigate to new order details page
      router.push(`/orders/${newOrder.id}`);
      router.refresh();
    } catch (err: unknown) {
      if (err instanceof Error) alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-white p-6 rounded-lg shadow"
    >
      <h2 className="text-xl font-semibold mb-4 text-gray-800">
        Add new Order
      </h2>

      <div>
        <label className="block text-sm font-bold text-gray-700 mb-2">
          Customer
        </label>
        <select
          value={customerID}
          onChange={(e) => setCustomerID(e.target.value)}
          className="w-full border border-gra-300 p-3 rounded-lg bg-gray-50 text-gray-900"
          required
        >
          <option value="">-- Select Customer --</option>
          {customers.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-bold text-gray-700 mb-2">
            Production Type
          </label>
          <select
            value={productionType}
            onChange={(e) => setProductionType(e.target.value)}
            className="w-full border border-gray-300 p-3 rounded-lg bg-gray-50 text-gray-900"
          >
            <option value="CMT">CMT (Cut, Make, Trim)</option>
            <option value="FOB">FOB (Full Package)</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-bold text-gray-700 mb-2">
            Total Quantity (pcs)
          </label>
          <input
            type="number"
            value={totalQuantity}
            onChange={(e) => setTotalQuantity(e.target.value)}
            className="w-full border border-gray-300 p-3 rounded-lg bg-gray-100 text-gray-900 cursor-not-allowed"
            placeholder="Auto-calculated from Size"
            readOnly
            required
          />
        </div>
      </div>

      <div className="border-t pt-6">
        <div className="flex justify-between items-center mb-4">
          <label className="block text-sm font-bold text-gray-700">
            Size Breakdown
          </label>
          <button
            type="button"
            onClick={addSizeRow}
            className="text-sm bg-gray-200 text-gray-800 px-3 py-1 rounded hover:bg-gray-300 font-medium"
          >
            + add Row
          </button>
        </div>

        <div className="space-y-3">
          {sizes.map((s, index) => (
            <div key={index} className="flex gap-3 items-center">
              <input
                type="text"
                value={s.size}
                onChange={(e) =>
                  handleSizeChange(index, "size", e.target.value)
                }
                placeholder="Size (e.g XL)"
                className="w-1/3 border border-gray-300 p-2 rounded bg-gray-50 text-gray-900"
              />
              <input
                type="number"
                value={s.quantity}
                onChange={(e) =>
                  handleSizeChange(index, "quantity", e.target.value)
                }
                placeholder="Qty"
                className="w-1/3 border border-gray-300 p-2 rounded bg-gray-50 text-gray-900"
              />
              <button
                type="button"
                onClick={() => removeSizeRow(index)}
                className="text-red-500 hover:text-red-700 font-bold px-2"
              >
                ✕
              </button>
            </div>
          ))}
        </div>
      </div>

      <div className="border-t pt-6">
        <label className="block text-sm font-bold text-gray-700 mb-2">
          Tech Packs & Sketches
        </label>
        <div className="border-2 border-dashed border-gray-400 rounded-lg p-6 flex flex-col items-center justify-center bg-gray-50">
          <input
            type="file"
            multiple
            accept="image/*"
            onChange={(e) => setFiles(e.target.files)}
            className="mb-2 text-sm text-gray-600 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
          />
          <p className="text-xs text-gray-500">
            Upload multiple sketches or reference images
          </p>
        </div>
      </div>

      <div className="pt-4">
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white font-bold py-3 px-4 rounded-lg hover:bg-blue-700 transition disabled:bg-gray-400"
        >
          {loading
            ? "Creating Order & Uploading..."
            : "Create Order & Start Sampling"}
        </button>
      </div>
    </form>
  );
}
