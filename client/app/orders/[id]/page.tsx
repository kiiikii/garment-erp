import Image from "next/image";

interface OrderSize {
  size_label: string;
  quantity: number;
}

interface OrderImage {
  id: number;
  image_url: string;
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
    phone: string;
    address: string;
  };
  sizes: OrderSize[] | null;
  images: OrderImage[] | null;
}

interface PageProps {
  params: Promise<{ id: string }>;
}

export default async function OrderDetailPage({ params }: PageProps) {
  const resolvedParams = await params;
  const orderID = resolvedParams.id;

  const res = await fetch(`http://localhost:8080/orders?id=${orderID}`, {
    cache: "no-cache",
  });

  if (!res.ok) {
    return (
      <main className="min-h-screen bg-gray-50 p-8 font-sans">
        <div className="max-w-3xl mx-auto bg-white p-6 rounded-lg shadow text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-2">
            Order Not Found
          </h1>
          <p className="text-gray-600 mb-4">
            Could not retrieve order #{orderID} from the backend.
          </p>
          <a
            href="/orders"
            className="text-blue-600 hover:underline font-medium"
          >
            ← Back to Orders Dashboard
          </a>
        </div>
      </main>
    );
  }

  const orders: Order = await res.json();

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Navigation & Header */}
        <div className="flex justify-between items-center">
          <div>
            <a
              href="/orders"
              className="text-blue-600 hover:underline text-sm font-medium"
            >
              ← Back to All Orders
            </a>
            <h1 className="text-3xl font-bold text-gray-800 mt-1">
              Order #{orders.id} Details
            </h1>
          </div>
          <span
            className={`px-3 py-1.5 rounded-full text-sm font-bold ${
              orders.status === "IN_PRODUCTION"
                ? "bg-green-100 text-green-800"
                : orders.status === "LOA_SIGNED"
                  ? "bg-purple-100 text-purple-800"
                  : orders.status === "SAMPLE_APPROVED"
                    ? "bg-blue-100 text-blue-800"
                    : "bg-yellow-100 text-yellow-800"
            }`}
          >
            {orders.status}
          </span>
        </div>

        {/* Grid Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Customer & Production Info */}
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
              Customer & Spec
            </h2>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Customer Name
              </p>
              <p className="text-gray-800 font-medium">
                {orders.customer?.name}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Production Type
              </p>
              <p className="text-blue-600 font-bold">{orders.production_type}</p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Total Quantity
              </p>
              <p className="text-gray-800 font-medium">
                {orders.total_quantity} pcs
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Assigned Layout / Sewing Line
              </p>
              <p className="font-semibold text-gray-700">
                {orders.layout_id ? (
                  `Sewing Line #${orders.layout_id}`
                ) : (
                  <span className="text-red-500 font-normal">
                    Unassigned (PPIC Gate Locked)
                  </span>
                )}
              </p>
            </div>
          </div>

          {/* Timelines */}
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">
              Sampling & Timelines
            </h2>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Internal Deadline (4 Days)
              </p>
              <p className="text-gray-800 font-medium">
                {new Date(orders.internal_sample_deadline).toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Customer Deadline (7 Days)
              </p>
              <p className="text-gray-800 font-medium">
                {new Date(orders.customer_sample_deadline).toLocaleDateString()}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase font-semibold">
                Actual Sample Finished At
              </p>
              <p className="font-medium text-gray-800">
                {orders.actual_sample_finished_at ? (
                  new Date(orders.actual_sample_finished_at).toLocaleString()
                ) : (
                  <span className="text-yellow-600">Pending completion</span>
                )}
              </p>
            </div>
          </div>
        </div>

        {/* Size Breakdown Card */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">
            Size Breakdown
          </h2>
          <div className="flex flex-wrap gap-3">
            {orders.sizes && orders.sizes.length > 0 ? (
              orders.sizes.map((s, idx) => (
                <div
                  key={idx}
                  className="bg-gray-100 border border-gray-200 rounded-lg p-3 text-center min-w-22.5"
                >
                  <p className="text-xs text-gray-500 font-bold uppercase">
                    Size {s.size_label}
                  </p>
                  <p className="text-lg font-bold text-gray-800">
                    {s.quantity}
                  </p>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm">
                No size quantities recorded.
              </p>
            )}
          </div>
        </div>

        {/* Tech Packs / Images Gallery */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">
            Tech Pack Images & Sketches
          </h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
            {orders.images && orders.images.length > 0 ? (
              orders.images.map((img) => (
                <div
                  key={img.id}
                  className="border rounded-lg overflow-hidden bg-gray-50"
                >
                  {/* Pointing directly to Go static server uploads */}
                  <Image
                    src={`http://localhost:8080/${img.image_url}`}
                    alt="Tech Pack"
                    className="w-full h-40 object-cover"
                  />
                  <div className="p-2 text-xs text-gray-500 text-center truncate">
                    Image ID: #{img.id}
                  </div>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm col-span-full">
                No tech pack images uploaded for this order yet.
              </p>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
