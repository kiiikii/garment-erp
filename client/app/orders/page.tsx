interface Order {
  id: number;
  quantity: number;
  production_type: string;
  created_at: string;
  customer: {
    id: number;
    name: string;
  };
}

export default async function OrdersPage() {
  const res = await fetch("http://localhost:8080/orders", {
    cache: "no-cache",
  });
  const orders: Order[] = await res.json();

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-800">
            Production Orders
          </h1>
          {/* A simple link to navigate back to the homepage */}
          <a href="/" className="text-blue-600 hover:underline font-medium">
            ← Back to Customers
          </a>
        </div>

        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full border-collapse">
            <thead className="bg-gray-800">
              <tr>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Order ID
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Customer Name
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Quantity
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Type
                </th>
              </tr>
            </thead>
            <tbody>
              {orders.map((o) => (
                <tr
                  key={o.id}
                  className="border-b border-gray-200 hover:bg-gray-100 transition-colors"
                >
                  <td className="py-3 px-4 text-gray-700 font-medium">
                    #{o.id}
                  </td>
                  {/* 3. Accessing the joined relational data! */}
                  <td className="py-3 px-4 text-gray-700">{o.customer.name}</td>
                  <td className="py-3 px-4 text-gray-700">{o.quantity}</td>
                  <td className="py-3 px-4 text-gray-700 font-bold">
                    {o.production_type}
                  </td>
                </tr>
              ))}

              {orders.length === 0 && (
                <tr>
                  <td colSpan={4} className="py-4 text-center text-gray-500">
                    No active orders found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </main>
  );
}
