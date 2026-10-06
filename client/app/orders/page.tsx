import Link from "next/link";

interface OrderSize {
  size_label: string;
  quantity: number;
}

interface Order {
  id: number;
  total_quantity: number;
  production_type: string;
  status: string;
  created_at: string;
  internal_sample_deadline: string;
  customer_sample_deadline: string;
  actual_sample_finished_at: string | null;
  layout_id: number | null;
  customer: {
    id: number;
    name: string;
  };
  sizes: OrderSize[] | null;
  internal_production_deadline: string;
  customer_production_deadline: string;
}

export default async function OrdersPage() {
  const res = await fetch("http://localhost:8080/orders", {
    cache: "no-cache",
  });
  const orders: Order[] = await res.json();

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-7xl px-6 mx-auto">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-800">
            Production Orders
          </h1>
          {/* A simple link to navigate back to the homepage */}
          <div>
            <Link
              href="/"
              className="text-blue-600 hover:underline font-medium mr-6"
            >
              ← Back to Dashboard
            </Link>
            <Link
              href="/orders/new"
              className="bg-blue-600 text-white px-4 py-2 rounded-lg font-bold hover:bg-blue-700 shadow"
            >
              + Create Order
            </Link>
          </div>
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
                  Type
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Quantity
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Sizes
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Status
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Internal Sample
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Cust Sample
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Internal Production
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Cust Production
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Layout ID
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
                    <a
                      href={`/orders/${o.id}`}
                      className="text-blue-600 hover:underline"
                    >
                      #{o.id}
                    </a>
                  </td>
                  {/* 3. Accessing the joined relational data! */}
                  <td className="py-3 px-4 text-gray-700">{o.customer.name}</td>
                  <td className="py-3 px-4 text-gray-700 font-bold">
                    {o.production_type}
                  </td>
                  <td className="py-3 px-4 text-gray-700">
                    {o.total_quantity}
                  </td>

                  {/* Size Breakdown Column */}
                  <td className="py-3 px-4 text-gray-600 text-xs">
                    {o.sizes ? (
                      o.sizes.map((s, idx) => (
                        <span
                          key={idx}
                          className="inline-block bg-gray-100 rounded px-1.5 py-0.5 mr-1 mb-1"
                        >
                          {s.size_label}: {s.quantity}
                        </span>
                      ))
                    ) : (
                      <span className="text-gray-400">No sizes</span>
                    )}
                  </td>

                  {/* Status Badge */}
                  <td className="py-3 px-4 ">
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-bold ${
                        o.status === "IN_PRODUCTION"
                          ? "bg-green-100 text-green-800"
                          : o.status === "LOA_SIGNED"
                            ? "bg-purple-100 text-purple-800"
                            : o.status === "SAMPLE_APPROVED"
                              ? "bg-blue-100 text-blue-800"
                              : "bg-yellow-100 text-yellow-800"
                      }`}
                    >
                      {o.status}
                    </span>
                  </td>

                  {/* Deadline */}
                  <td className="py-3 px-4 text-gray-600 text-xs">
                    {new Date(o.internal_sample_deadline).toLocaleDateString()}
                  </td>
                  <td className="py-3 px-4 text-gray-600 text-xs">
                    {new Date(o.customer_sample_deadline).toLocaleDateString()}
                  </td>
                  <td className="py-3 px-4 text-gray-600 text-xs">
                    {new Date(
                      o.internal_production_deadline,
                    ).toLocaleDateString()}
                  </td>
                  <td className="py-3 px-4 text-gray-600 text-xs">
                    {new Date(
                      o.customer_production_deadline,
                    ).toLocaleDateString()}
                  </td>

                  {/* Layout Assignment */}
                  <td className="py-3 px-4 font-semibold text-gray-700">
                    {o.layout_id ? (
                      `Line #${o.layout_id}`
                    ) : (
                      <span className="bg-red-500 px-2 py-1 rounded-full text-white text-xs font-bold">
                        Unassigned
                      </span>
                    )}
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
