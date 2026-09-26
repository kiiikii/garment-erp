import CustomerForm from "@/components/CustomerForm";

//! Define typescript Interface (match with Go Customer Response struct)
interface Customer {
  id: number;
  name: string;
  phone: string;
  address: string;
}

//! Next.js server component
export default async function Home() {
  //! fetch data directly from backend
  //! the {cache: "no-store"} ensure Next.js always gets a live data from database
  const res = await fetch("http://localhost:8080/customers", {
    cache: "no-cache",
  });

  //! parsing JSON respon into typescript array
  const customers: Customer[] = await res.json();

  //! render the HTML and inject the dynamic data
  return (
    <main className="min-h-screen bg-gray-50 p-8 font-sans">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6 text-gray-800">
          Garment ERP Dashboard
        </h1>

        <CustomerForm />

        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full border-collapse">
            <thead className="bg-gray-800">
              <tr>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  ID
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Name
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Phone
                </th>
                <th className="py-3 px-4 text-left text-sm font-semibold text-white">
                  Address
                </th>
              </tr>
            </thead>
            <tbody>
              {/* Loop through the customers array and create a row for each one */}
              {customers.map((c) => (
                <tr
                  key={c.id}
                  className="border-b border-gray-200 hover:bg-gray-100 transition-colors"
                >
                  <td className="py-3 px-4 text-gray-700">{c.id}</td>
                  <td className="py-3 px-4 text-gray-700 font-medium">
                    {c.name}
                  </td>
                  <td className="py-3 px-4 text-gray-700">{c.phone}</td>
                  <td className="py-3 px-4 text-gray-700">{c.address}</td>
                </tr>
              ))}

              {/* Fallback if the database is empty */}
              {customers.length === 0 && (
                <tr>
                  <td colSpan={4} className="py-4 text-center text-gray-500">
                    No customers found in the database.
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
