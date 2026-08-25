import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import type { Key } from "react-aria-components";

export type Duration = string;

interface ExpiryDurationSelectProps {
  onSelectionChange: (duration: Duration) => void;
};

const ExpiryDurationSelect = (props: ExpiryDurationSelectProps) => {
  const durations = [
    // TODO use Temporal to automatically generate labels from duration
    {
      label: "10 min",
      duration: "PT10M",
    },
    {
      label: "1 hour",
      duration: "PT1H",
    },
    {
      label: "1 day",
      duration: "P1D",
    },
    {
      label: "7 days",
      duration: "P7D",
    },
  ];

  const { onSelectionChange } = props;

  const handleSelectionChange = (keys: Set<Key>) => {
    const selectedValue = keys.keys().next().value;

    if (keys.size != 1 || selectedValue === undefined) {
      console.error(
        `Expected ExpiryDurationSelect to only have 1 selection, but has ${keys.size.toString()}`,
      );
      return;
    }

    onSelectionChange(selectedValue.toString());
  };

  return (
    <ToggleGroup
      defaultSelectedKeys={[durations[0].duration]}
      variant="outline"
      onSelectionChange={handleSelectionChange}
      disallowEmptySelection={true}
    >
      {durations.map((duration) => (
        <ToggleGroupItem key={duration.duration} id={duration.duration}>
          {duration.label}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  );
};

export default ExpiryDurationSelect;
