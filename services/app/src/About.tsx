// import React from 'react';

// function About() {
//   return (
//     <div className="About">
//       <p>Hello</p>
//       <header className="About-header">
//       </header>
//     </div>
//   );
// }
import React, { useState } from 'react';

interface TodoItem {
  id: number,
  title: string,
}

function About() {
  const [data, setData] = useState<TodoItem[]>([]);
  
  const fetchData = () => {
    fetch('https://jsonplaceholder.typicode.com/todos')
      .then(response => response.json())
      .then(data => setData(data))
      .catch(error => console.error('Error fetching data:', error));
  };

  return (
    <div>
      <button onClick={fetchData}>Fetch Data</button>
      <ul>
        {data.map(item => (
          <li> {item.title} </li>
        ))}
      </ul>
    </div>
  );
}

export default About;