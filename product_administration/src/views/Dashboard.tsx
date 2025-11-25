export default function Dashboard() {
    return (
        <div>
            <h1 className="text-2xl font-bold text-slate-900 mb-4">Irányítópult</h1>
            <p className="text-gray-500">Üdvözöljük az adminisztrációs felületen.</p>
            <button className="px-6 py-3 mt-4 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all shadow-sm shadow-blue-200 focus:outline-none focus:ring-4 focus:ring-blue-100 disabled:opacity-50 disabled:cursor-not-allowed">Ugrás a központi felületre</button>
        </div>
    );
}