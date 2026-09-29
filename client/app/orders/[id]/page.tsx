import Image from "next/image"

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

export default async function OrderDetailPage({ params }: { params: { id: string } }) {
  const resolvedParams = await params
  const id = resolvedParams.id

  const res = await fetch(`http://localhost:8080/orders/${id}`, {
    cache: "no-cache",
  });
  
  if (!res.ok) {
    return <div className="p-8 text-red-600 font-bold">Order not found or backend error.</div>;
  }

  const data = await res.json();
  
  // Handle if backend returns an array or a single object
  const order: Order = Array.isArray(data) ? data[0] : data;

  if (!order) {
    return <div className="p-8 text-red-600 font-bold">Order #{id} not found in database.</div>;
  }

  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-6">
        
        {/* Navigation & Header */}
        <div className="flex justify-between items-center">
          <div>
            <a href="/orders" className="text-blue-600 hover:underline text-sm font-medium">
              ← Back to Orders Dashboard
            </a>
            <h1 className="text-3xl font-bold text-gray-900 mt-1">
              Order #{order.id} — {order.customer?.name || "Unknown Customer"}
            </h1>
          </div>
          <span className={`px-3 py-1 rounded-full text-sm font-bold ${
            order.status === 'IN_PRODUCTION' ? 'bg-green-100 text-green-800' :
            order.status === 'LOA_SIGNED' ? 'bg-purple-100 text-purple-800' :
            'bg-yellow-100 text-yellow-800'
          }`}>
            {order.status || "PENDING"}
          </span>
        </div>

        {/* Timelines & Factory Specs */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">Production Specs</h2>
            <p><strong className="text-gray-600">Production Type:</strong> <span className="text-blue-600 font-semibold">{order.production_type}</span></p>
            <p className="text-gray-900"><strong className="text-gray-600">Total Quantity:</strong> {order.total_quantity} pcs</p>
            <p><strong className="text-gray-600">Assigned Layout:</strong> {order.layout_id ? `Sewing Line #${order.layout_id}` : <span className="text-red-500 font-medium">Unassigned (PPIC Gate Locked)</span>}</p>
          </div>

          <div className="bg-white p-6 rounded-lg shadow space-y-4">
            <h2 className="text-lg font-bold text-gray-800 border-b pb-2">Sampling Timelines</h2>
            <p className="text-gray-900"><strong className="text-gray-600">Internal Deadline (4 Days):</strong> {order.internal_sample_deadline ? new Date(order.internal_sample_deadline).toLocaleDateString() : "N/A"}</p>
            <p className="text-gray-900"><strong className="text-gray-600">Customer Deadline (7 Days):</strong> {order.customer_sample_deadline ? new Date(order.customer_sample_deadline).toLocaleDateString() : "N/A"}</p>
            <p>
              <strong className="text-gray-600">Actual Finished At:</strong>{" "}
              {order.actual_sample_finished_at ? (
                <span className="text-green-600 font-semibold">{new Date(order.actual_sample_finished_at).toLocaleString()}</span>
              ) : (
                <span className="text-amber-500 font-medium">Pending Sampling Completion</span>
              )}
            </p>
          </div>
        </div>

        {/* Size Breakdown */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">Size Breakdown</h2>
          <div className="flex gap-4">
            {order.sizes && order.sizes.length > 0 ? (
              order.sizes.map((s, idx) => (
                <div key={idx} className="bg-gray-100 px-4 py-2 rounded text-center border">
                  <div className="text-xs text-gray-500 font-bold uppercase">{s.size_label}</div>
                  <div className="text-lg font-bold text-gray-800">{s.quantity}</div>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm italic">No sizes specified.</p>
            )}
          </div>
        </div>

        {/* Uploaded Tech Packs / Images */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-lg font-bold text-gray-800 border-b pb-2 mb-4">Tech Packs & Design Sketches</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {order.images && order.images.length > 0 ? (
              order.images.map((img) => (
                <div key={img.id} className="border rounded p-2 bg-gray-50 relative h-32">
                  <Image 
                    src={`http://localhost:8080/${img.image_url}`} 
                    alt="Tech Pack" 
                    fill
                    sizes="(max-width: 768px) 50vw, 25vw"
                    className="object-cover rounded"
                  />
                  <span className="text-xs text-gray-500 mt-1 block truncate relative z-10 bg-white/80 px-1 rounded">Image #{img.id}</span>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm italic">No tech packs uploaded for this order yet.</p>
            )}
          </div>
        </div>

      </div>
    </main>
  );
}
