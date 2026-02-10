interface TabsProps {
  selected: string;
  onSelect: (tab: string) => void;
  items: string[];
}

export function Tabs({ selected, onSelect, items }: TabsProps) {
  return (
    <nav className="tabs" aria-label="section tabs">
      {items.map((item) => (
        <button
          key={item}
          className={item === selected ? 'tab active' : 'tab'}
          onClick={() => onSelect(item)}
          type="button"
        >
          {item}
        </button>
      ))}
    </nav>
  );
}
