import React, { useState } from 'react';

type Release = {
    style_genre: string
    artist_name: string
    release_count: number
}

const SearchComponent: React.FC = () => {
    // States to hold dropdown selection and input value
    const [searchType, setSearchType] = useState<'artists' | 'styles'>('artists');
    const [inputValue, setInputValue] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [releases, setReleases] = useState<Release[] | null>(null);
    const [emptyResult, setEmptyResult] = useState<boolean | null>(null)

    // Handler for dropdown selection
    const handleSearchTypeChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        setSearchType(event.target.value as 'artists' | 'styles');
    };

    // Handler for text input change
    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(event.target.value);
    };

    // Function to make the HTTP request
    const handleSearch = async () => {
        if (!inputValue) {
            setError('Please enter a search term');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const response = await fetch(`http://localhost:8081/${searchType}?${searchType.substring(0, searchType.length -1)}=${inputValue}&sort=release_count&order=desc`);
            if (!response.ok) {
                throw new Error('Failed to fetch data');
            }
            const data = await response.json();
            console.log('Search results:', data);
            setReleases(data)
            if (data == null) {
                setEmptyResult(true)
            } else {
                setEmptyResult(false)
            }
        } catch (error) {
            // @ts-ignore
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <label>
                Search by:
                <select value={searchType} onChange={handleSearchTypeChange}>
                    <option value="artists">Artist</option>
                    <option value="styles">Style/Genre</option>
                </select>
            </label>

            <input
                type="text"
                value={inputValue}
                onChange={handleInputChange}
                placeholder="Enter search term"
            />

            <button onClick={handleSearch} disabled={loading}>
                {loading ? 'Searching...' : 'Search'}
            </button>

            {error && <p style={{ color: 'red' }}>{error}</p>}
            {/* Render the table if data is available */}
            {releases && (
                <table style={{marginLeft:"auto", marginRight:"auto"}}>
                    <thead>
                    <tr>
                        <th>Style/Genre</th>
                        <th>Artist Name</th>
                        <th>Release Count</th>
                    </tr>
                    </thead>
                    <tbody>
                    {releases.map((release, index) => (
                        <tr key={index}>
                            <td>{release.style_genre}</td>
                            <td>{release.artist_name}</td>
                            <td>{release.release_count}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
            {emptyResult && emptyResult? <div>Result was empty</div>:null}
        </div>
    );
};

export default SearchComponent;
